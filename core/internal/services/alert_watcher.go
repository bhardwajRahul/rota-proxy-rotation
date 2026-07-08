package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alpkeskin/rota/core/internal/models"
	"github.com/alpkeskin/rota/core/internal/repository"
	"github.com/alpkeskin/rota/core/pkg/logger"
)

// AlertWatcher monitors pool health and fires webhook alerts when active proxies
// drop below the configured threshold.
type AlertWatcher struct {
	poolRepo *repository.PoolRepository
	log      *logger.Logger
	client   *http.Client
	interval time.Duration
}

// NewAlertWatcher creates a new AlertWatcher with a 2-minute check interval.
func NewAlertWatcher(poolRepo *repository.PoolRepository, log *logger.Logger) *AlertWatcher {
	return &AlertWatcher{
		poolRepo: poolRepo,
		log:      log,
		client:   &http.Client{Timeout: 10 * time.Second},
		interval: 2 * time.Minute,
	}
}

// Start begins the background watcher loop.
func (w *AlertWatcher) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.check(ctx)
			}
		}
	}()
	w.log.Info("alert watcher started", "interval", w.interval)
}

// check loads all enabled rules and fires those whose thresholds are exceeded
// and whose cooldown has elapsed.
func (w *AlertWatcher) check(ctx context.Context) {
	rules, err := w.poolRepo.GetAllAlertRules(ctx)
	if err != nil {
		w.log.Warn("alert watcher: failed to load rules", "error", err)
		return
	}
	if len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		pool, err := w.poolRepo.GetByID(ctx, rule.PoolID)
		if err != nil || pool == nil {
			continue
		}

		if pool.ActiveProxies >= rule.MinActiveProxies {
			continue // threshold OK
		}

		// Check cooldown
		if rule.LastFiredAt != nil {
			cooldown := time.Duration(rule.CooldownMinutes) * time.Minute
			if time.Since(*rule.LastFiredAt) < cooldown {
				continue // still in cooldown
			}
		}

		w.log.Warn("pool alert threshold triggered",
			"pool_id", pool.ID,
			"pool_name", pool.Name,
			"active", pool.ActiveProxies,
			"threshold", rule.MinActiveProxies,
		)

		if err := w.fire(ctx, rule, *pool); err != nil {
			w.log.Error("failed to fire pool alert webhook",
				"rule_id", rule.ID,
				"url", redactWebhookURL(rule.WebhookURL),
				"error", err,
			)
		} else {
			// Record fire time
			if err := w.poolRepo.UpdateAlertRuleFiredAt(ctx, rule.ID); err != nil {
				w.log.Warn("failed to update alert rule fired_at", "rule_id", rule.ID, "error", err)
			}
		}
	}
}

// redactWebhookURL returns a form of the webhook URL that is safe to log. For
// Telegram (and any similar /bot<token>/ path), the token segment embeds a
// secret that the logger's DB hook would otherwise persist — so it is replaced
// with a placeholder. Query strings (which may carry chat_id etc.) are dropped.
func redactWebhookURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "[redacted]"
	}
	path := u.Path
	if strings.HasPrefix(path, "/bot") {
		rest := strings.TrimPrefix(path, "/bot")
		if i := strings.IndexByte(rest, '/'); i != -1 {
			path = "/bot<redacted>" + rest[i:]
		} else {
			path = "/bot<redacted>"
		}
	}
	return u.Scheme + "://" + u.Host + path
}

// fire sends the webhook request for a triggered alert rule.
//
// Telegram Bot API endpoints are detected by host and translated into a
// sendMessage call, since Telegram won't accept the generic JSON payload
// (issue #31). For a group topic, configure the webhook URL as:
//
//	https://api.telegram.org/bot<TOKEN>/sendMessage?chat_id=<ID>&message_thread_id=<TOPIC_ID>
//
// Any other host receives the generic PoolAlertPayload JSON as before.
func (w *AlertWatcher) fire(ctx context.Context, rule models.PoolAlertRule, pool models.ProxyPool) error {
	if u, err := url.Parse(rule.WebhookURL); err == nil && strings.EqualFold(u.Hostname(), "api.telegram.org") {
		return w.fireTelegram(ctx, rule, pool, u)
	}

	payload := models.PoolAlertPayload{
		Event:         "pool.degraded",
		PoolID:        pool.ID,
		PoolName:      pool.Name,
		ActiveProxies: pool.ActiveProxies,
		TotalProxies:  pool.TotalProxies,
		Threshold:     rule.MinActiveProxies,
		FiredAt:       time.Now().UTC(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	method := rule.WebhookMethod
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequestWithContext(ctx, method, rule.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Rota-AlertWatcher/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned HTTP %d", resp.StatusCode)
	}

	w.log.Info("pool alert webhook fired",
		"rule_id", rule.ID,
		"pool_id", pool.ID,
		"status", resp.StatusCode,
	)
	return nil
}

// fireTelegram sends the alert as a Telegram Bot API sendMessage call.
// chat_id (required) and optional message_thread_id (for group topics) are read
// from the configured webhook URL's query string; the bot token stays in the
// URL path. See fire() for the expected URL format (issue #31).
func (w *AlertWatcher) fireTelegram(ctx context.Context, rule models.PoolAlertRule, pool models.ProxyPool, u *url.URL) error {
	q := u.Query()
	chatID := q.Get("chat_id")
	if chatID == "" {
		return fmt.Errorf("telegram webhook missing chat_id query parameter")
	}

	text := fmt.Sprintf(
		"🔴 Rota pool alert\nPool: %s (#%d)\nActive proxies: %d / %d\nThreshold: %d",
		pool.Name, pool.ID, pool.ActiveProxies, pool.TotalProxies, rule.MinActiveProxies,
	)

	msg := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	if tid := q.Get("message_thread_id"); tid != "" {
		msg["message_thread_id"] = tid
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal telegram message: %w", err)
	}

	// Telegram requires the bot token path segment to be prefixed with "bot"
	// (e.g. /bot<TOKEN>/sendMessage). Users commonly paste the token without it,
	// which makes the API return 404 — auto-correct it so the alert still works.
	path := u.Path
	if !strings.HasPrefix(path, "/bot") {
		path = "/bot" + strings.TrimPrefix(path, "/")
	}

	// Send to the bare endpoint (scheme+host+path) — query params carried the
	// routing info we've now folded into the JSON body.
	endpoint := u.Scheme + "://" + u.Host + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram returned HTTP %d", resp.StatusCode)
	}

	w.log.Info("pool alert telegram fired",
		"rule_id", rule.ID,
		"pool_id", pool.ID,
		"status", resp.StatusCode,
	)
	return nil
}

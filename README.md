<div align="center" style="margin-bottom: 20px;">
  <img src="static/rota_logo.png" alt="rota" width="100px">
  <h1 align="center">
  Rota - Proxy Rotation Platform
  </h1>
</div>

<p align="center">
<a href="https://opensource.org/licenses/Apache-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.25.3-00ADD8?logo=go"></a>
<a href="https://nextjs.org"><img src="https://img.shields.io/badge/Next.js-16-000000?logo=next.js"></a>
<a href="https://www.timescale.com/"><img src="https://img.shields.io/badge/TimescaleDB-2.22-FDB515?logo=timescale"></a>
<a href="https://github.com/alpkeskin/rota/releases"><img src="https://img.shields.io/github/release/alpkeskin/rota"></a>
<a href="https://github.com/alpkeskin/rota/actions"><img src="https://img.shields.io/github/actions/workflow/status/alpkeskin/rota/release.yaml"></a>
</p>


![Khipu Screenshot](static/dashboard.png)


## 🎯 Overview

**Rota** is a modern, full-stack proxy rotation platform that combines enterprise-grade proxy management with a beautiful, real-time web dashboard. Built with performance and scalability in mind, Rota handles thousands of requests per second while providing comprehensive monitoring, analytics, and control through an intuitive interface.

Whether you're conducting web scraping operations, performing security research, load testing, or need reliable proxy management at scale, Rota delivers a complete solution with:

- **High-Performance Core**: Lightning-fast Go-based proxy server with intelligent rotation strategies
- **Real-Time Dashboard**: Modern Next.js web interface with live metrics and monitoring
- **Time-Series Analytics**: TimescaleDB-powered storage for historical analysis and insights
- **Production-Ready**: Docker-based deployment with health checks, graceful shutdown, and monitoring

---

## ✨ Key Features

### Core Proxy Server
- 🚀 **High Performance**: Handle thousands of concurrent requests with minimal latency — pooled upstream transports (keep-alive reuse), zero-copy `splice(2)` tunneling on Linux, and batched request telemetry that coalesces per-request DB writes
- 🔄 **Smart Rotation**: Multiple rotation strategies (random, round-robin, least connections, time-based)
- 🤖 **Automatic Management**: Real-time proxy pool monitoring with automatic unhealthy proxy removal
- 🌍 **Multi-Protocol**: Full support for HTTP, HTTPS, SOCKS4, SOCKS4A, and SOCKS5
- ✅ **Health Checking**: Built-in proxy validation to maintain a healthy pool
- 🔒 **Authentication**: Basic auth support for proxy server
- ⚡ **Rate Limiting**: Configurable rate limiting to prevent abuse
- 🔗 **Proxy Chaining**: Compatible with upstream proxies (Burp Suite, OWASP ZAP, etc.)
- ⏱️ **Configurable Timeouts**: Fine-grained control over request timeouts and retries
- 🔁 **Redirect Support**: Optional HTTP redirect following

### Proxy Sources & Auto-Import
- 📥 **Remote TXT Lists**: Add URLs pointing to `ip:port` proxy lists — fetched automatically on schedule
- 🕐 **Per-Source Interval**: Each source has its own refresh interval (in minutes)
- 🔁 **Background Scheduler**: Overdue sources are fetched automatically every minute
- 🌍 **Protocol per Source**: Assign HTTP, HTTPS, SOCKS4, SOCKS4a, or SOCKS5 to each list

### GeoIP & Geo Distribution
- 🗺️ **Automatic GeoIP**: Proxies are geolocated via [ip-api.com](http://ip-api.com) (free, no API key required)
- 🏙️ **City-Level Data**: Country, region, city, ISP, latitude, longitude per proxy
- 🔍 **Geo Explorer**: Expandable country tree with city drill-down in the dashboard
- ♻️ **Auto-Enrich**: Geo data updated automatically after every source fetch

### Proxy Pools
- 🗂️ **Named Pools**: Group proxies by any combination of countries, cities, ISPs, or custom tags
- ☑️ **Multi-Filter Builder**: Pick geo locations, ISP substrings, or proxy tags — mix freely in one pool
- 🔄 **Auto / Manual Sync**: `sync_mode: auto` rebuilds membership on every import; `manual` keeps it frozen until you trigger sync explicitly
- 🔁 **Rotation Strategies**: Per-pool `roundrobin`, `random`, or `sticky` (hold N requests per IP)
- ⚡ **Async Health Checks**: Run health checks against any URL; progress shown in real time
- ⏱️ **Scheduled Checks**: Cron-style schedule per pool (`*/30 * * * *`)
- 📤 **Export**: Download pool proxy list as `.txt` or `.csv` (`GET /api/v1/pools/{id}/export?format=txt|csv`)
- 🔔 **Webhook Alerts**: Per-pool alert rules — fire a POST/GET webhook when active proxy count drops below threshold, with configurable cooldown

### Per-User Pool Authentication
- 👤 **Proxy Users**: Create users with bcrypt passwords, each assigned a main pool + ordered fallbacks
- 🔗 **Usage**: `http://user:pass@host:8000` — the proxy routes through the user's pool chain
- 🔄 **Automatic Failover**: If a pool has no live IPs, requests cascade to fallback pools
- 🔁 **Retry Logic**: Each retry picks a fresh proxy; failed IPs are excluded for that request
- 📊 **Full Tracking**: All requests, success rates, and response times tracked per proxy
- ⚡ **Per-User Rate Limit**: Optional `requests_per_minute` cap per user (0 = unlimited)

### Security
- 🔐 **JWT Authentication**: All API endpoints require a valid JWT token; the browser auto-redirects to login on expiry with "Session expired" message
- 🔑 **Bcrypt Admin Credentials**: Dashboard password stored as bcrypt hash in database
- 🔄 **Change Password**: Update username/password via the Settings UI (requires current password)
- 🌐 **Public endpoints only**: `GET /health` and `POST /auth/login`
- 🛡️ **Auth Brute-Force Protection**: Per-IP block after N failed attempts + global lockout when request rate exceeds threshold (all configurable via `.env`)
- 🧱 **Spoof-Resistant Rate Limiting**: Client IP is only derived from forwarded headers when `TRUST_PROXY_HEADERS` is enabled, so per-IP limits can't be bypassed with a forged `X-Forwarded-For`
- 🔌 **WebSocket Origin Validation**: Live dashboard/log streams reject cross-origin connections (CSWSH protection), honoring the configured CORS allowlist
- 🏷️ **Proxy Tags**: Label proxies with custom tags for fine-grained pool filtering
- 🧹 **Dead Proxy Cleanup**: Configurable automatic removal of long-failed or low-quality proxies

### Web Dashboard
- 📊 **Real-Time Metrics**: Live statistics, charts, and system monitoring
- 🔄 **Proxy Management**: Add, edit, delete, and test proxies through the UI
- 📝 **Live Logs**: WebSocket-based real-time log streaming
- 💻 **System Monitoring**: CPU, memory, disk, and runtime metrics
- ⚙️ **Configuration**: Manage settings through the web interface
- 🎨 **Modern UI**: Beautiful, responsive design with dark mode support
- 📱 **Mobile-Friendly**: Fully responsive across all devices

### Data & Analytics
- 📈 **Time-Series Storage**: TimescaleDB for efficient historical data storage
- 🔍 **Request History**: Track all proxy requests with detailed metadata
- 📉 **Performance Analytics**: Analyze proxy performance over time
- 🎯 **Usage Insights**: Understand traffic patterns and proxy utilization

### DevOps & Deployment
- 🐳 **Docker-Native**: Production-ready containerized deployment
- 🔧 **Easy Configuration**: All config via `.env` — see `.env.example` for all options
- 🏥 **Health Checks**: Built-in health endpoints for monitoring
- 🛑 **Graceful Shutdown**: Clean shutdown with connection draining
- 📊 **Observability**: Structured JSON logging and metrics endpoints

---

## 🚀 Quick Start

### Using Docker Compose (Recommended)

Everything runs behind a single entry point, so there's just one URL to open and
no API URL to configure.

```bash
# 1. Clone and start — no config file needed
git clone https://github.com/alpkeskin/rota.git
cd rota
docker compose up -d          # or: make up

# 2. Grab the auto-generated admin password from the logs
make password                  # or: docker compose logs rota-core | grep -i password
```

Then open **http://localhost** and log in with user `admin` and the password
from the logs. That's it — the dashboard, API and live logs are all served from
the same origin, so it works the same whether you're on `localhost` or a remote
server's IP.

**What's exposed:**
- 🌐 **Web UI + API**: http://localhost (everything — `/`, `/api`, `/docs`)
- 🔄 **Proxy**: `localhost:8000` (what your clients connect through)

> First-boot credentials are seeded once. Leave `ROTA_ADMIN_PASSWORD` unset to
> get a strong random password (shown in the logs), or set it in `.env` to pick
> your own. Change it anytime via **Settings → Admin Account**.

### Configuration

No `.env` is required — defaults work out of the box. Copy `cp .env.example .env`
only to change something. The common knobs:

| Variable | Default | Description |
|---|---|---|
| `SITE_ADDRESS` | `:80` | Web entry address. Set a domain for automatic HTTPS |
| `ROTA_ADMIN_PASSWORD` | _(random)_ | Initial admin password; blank → generated & logged |
| `ROTA_ADMIN_USER` | `admin` | Initial dashboard username (seeded once) |
| `PROXY_PORT` | `8000` | Host port for the proxy your clients connect to |
| `HTTP_PORT` / `HTTPS_PORT` | `80` / `443` | Web entry ports (Caddy) |
| `DB_PASSWORD` | `rota_password` | TimescaleDB password |
| `CORS_ALLOWED_ORIGINS` | `*` | API CORS allowlist (irrelevant behind the proxy) |
| `TRUST_PROXY_HEADERS` | `true` | Trust `X-Forwarded-For`/`X-Real-IP` for the login rate limiter. Keep `true` behind the bundled Caddy; set `false` if the API is exposed directly |
| `LOG_LEVEL` | `info` | Log verbosity: `debug`, `info`, `warn`, `error` |

See `.env.example` for the full list including auth brute-force protection.

> **Note**: `ROTA_ADMIN_USER` / `ROTA_ADMIN_PASSWORD` are only used when the
> database is empty (first start). Afterwards, use **Settings → Admin Account**.

### Production Deployment (HTTPS)

Point a domain at the server and set one variable — Caddy obtains and renews a
TLS certificate automatically:

```bash
# .env
SITE_ADDRESS=rota.example.com
DB_PASSWORD=a-strong-random-password
ROTA_ADMIN_PASSWORD=a-strong-password
```

```bash
docker compose up -d --build
```

Everything is then served over HTTPS at `https://rota.example.com` — no separate
API host, no dashboard rebuild when the domain changes.

### Using Docker

Pull and run the core service:

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/alpkeskin/rota:latest

# Run with basic configuration
docker run -d \
  --name rota-core \
  -p 8000:8000 \
  -p 8001:8001 \
  -e DB_HOST=your-db-host \
  -e DB_USER=rota \
  -e DB_PASSWORD=your-password \
  ghcr.io/alpkeskin/rota:latest
```

### From Source

```bash
# Prerequisites: Go 1.25.3+, Node.js 20+, pnpm, and TimescaleDB reachable

# Clone the repository
git clone https://github.com/alpkeskin/rota.git
cd rota

# Start the Go core (serves API on :8001, proxy on :8000)
make dev-core            # or: cd core && go run ./cmd/server

# In a second terminal, start the dashboard.
# Dev runs on separate ports, so point the browser at the core directly:
cd dashboard
cp .env.local.example .env.local     # sets NEXT_PUBLIC_API_URL=http://localhost:8001
make -C .. dev-dashboard             # or: pnpm install && pnpm dev
```

> The DB connection and admin credentials come from the same environment
> variables as the Docker setup (see `.env.example`).

### Testing the Proxy

```bash
# Route traffic through Rota proxy
curl -x http://localhost:8000 https://api.ipify.org?format=json

# Per-user pool routing (after creating a Proxy User in the dashboard)
curl -x http://myuser:mypassword@localhost:8000 https://api.ipify.org?format=json

# Using environment variables
export HTTP_PROXY=http://localhost:8000
export HTTPS_PROXY=http://localhost:8000
curl https://api.ipify.org?format=json
```

---

## 📚 API Documentation

### Interactive API Documentation (Swagger)

Rota provides interactive API documentation. Once the stack is running, you can access it at:

```
http://localhost/docs
```

The docs interface allows you to:
- 📖 Browse all available API endpoints
- 🧪 Test API requests directly from your browser
- 📝 View request/response schemas
- 🔍 Explore authentication requirements

**Quick Access:**
- **API docs**: http://localhost/docs
- **OpenAPI Spec**: http://localhost/api/v1/swagger.json

---

## 🏗️ Architecture

Rota is a monorepo. A single reverse proxy (Caddy) is the only web entry point,
so the browser talks to one origin; the dashboard, API and WebSockets are all
same-origin behind it. Only the proxy port is exposed separately.

```
                         Browser
                            │  http(s)://localhost  (one origin)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                        Rota Platform                        │
├─────────────────────────────────────────────────────────────┤
│                    ┌──────────────┐                         │
│                    │    Caddy     │  :80 / :443 (auto-HTTPS)│
│                    │ reverse proxy│                         │
│                    └──────┬───────┘                         │
│           /  , /_next     │      /api/* , /ws/* , /docs     │
│        ┌──────────────────┴───────────────┐                 │
│        ▼                                   ▼                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │   Dashboard  │    │  Core (API)  │───▶│ TimescaleDB  │   │
│  │   Next.js    │    │     Go       │    │  PostgreSQL  │   │
│  │  (internal)  │    │  (internal)  │    │  (internal)  │   │
│  └──────────────┘    └──────┬───────┘    └──────────────┘   │
│                             ▼                               │
│                      ┌──────────────┐                       │
│                      │ Proxy Server │  :8000 (exposed)      │
│                      │      Go      │                       │
│                      └──────┬───────┘                       │
└─────────────────────────────┼───────────────────────────────┘
                              ▼
                     ┌──────────────────┐
                     │   Proxy Pool     │
                     │  (External IPs)  │
                     └──────────────────┘
```

---

### Rotation Strategies

- **Random**: Select a random proxy for each request
- **Round Robin**: Distribute requests evenly across all proxies
- **Least Connections**: Route to the proxy with fewest active connections
- **Time-Based**: Rotate proxies at fixed intervals

---

## 🐳 Deployment

### Production Deployment

#### Using Docker Compose

```bash
# Production configuration
docker compose -f docker-compose.yml up -d

# Enable auto-restart
docker compose up -d --restart=unless-stopped
```

---

## 🗂️ Proxy Sources & Pools

### How Proxy Sources work

1. Go to **Proxy Sources** in the dashboard
2. Add a URL pointing to a plain-text proxy list (one `ip:port` per line)
3. Choose the protocol and refresh interval
4. Click **Fetch Now** or wait for the scheduler

The system will:
- Download and parse the list
- Upsert proxies into the database (duplicates ignored)
- Automatically look up GeoIP data for every new proxy
- Re-sync all pools that have `Auto-sync` enabled

### Geo Distribution & Pools

After proxies are geolocated, open the **Proxy Pools → Geo Distribution** tab:

- Browse all proxy-holding countries; click a country to expand cities
- Check individual countries or cities; mix them freely
- Click **Create Pool from selection** — the pool is created and filled instantly

Pools also support **ISP filters** (substring match, OR logic) and **tag filters** (AND logic — proxy must carry all specified tags). Combine geo + ISP + tags in any combination.

#### Pool Sync Modes

| Mode | Behaviour |
|------|-----------|
| `auto` | Pool membership is rebuilt automatically after every proxy import or geo-enrichment |
| `manual` | Membership only changes when you press **Sync** — useful for curated pools |

#### Exporting a Pool

```bash
# Plain text — one protocol://ip:port per line
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost/api/v1/pools/{id}/export?format=txt" -o pool.txt

# CSV — with status, geo, ISP, success rate
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost/api/v1/pools/{id}/export?format=csv" -o pool.csv
```

#### Webhook Alerts

Add an alert rule to a pool to be notified when the active proxy count drops below a threshold:

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  "http://localhost/api/v1/pools/{id}/alert-rules" \
  -d '{
    "enabled": true,
    "min_active_proxies": 10,
    "webhook_url": "https://hooks.slack.com/...",
    "cooldown_minutes": 30
  }'
```

Payload sent to the webhook:
```json
{
  "event": "pool.degraded",
  "pool_id": 1,
  "pool_name": "US Residential",
  "active_proxies": 3,
  "total_proxies": 50,
  "threshold": 10,
  "fired_at": "2026-04-02T04:30:00Z"
}
```

##### Telegram alerts (group topics)

Any other URL receives the generic JSON payload above. When the webhook host is
`api.telegram.org`, Rota instead calls the Telegram Bot API `sendMessage` method
and **generates the message text for you** — so you don't need a `text=`
parameter. To alert a topic inside a group:

1. Create a bot with [@BotFather](https://t.me/BotFather) and copy its token.
2. Add the bot to your group and give it permission to post in the topic.
3. Find the group's numeric `chat_id` (a negative number, e.g. `-1001234567890`)
   and the topic's `message_thread_id`.
4. Set the alert rule's **Webhook URL** to:

```
https://api.telegram.org/bot<TOKEN>/sendMessage?chat_id=<-100...>&message_thread_id=<topic_id>
```

- `message_thread_id` is optional — omit it to post to the group's main channel.
- The `bot` prefix on the token is required by Telegram; Rota adds it
  automatically if you leave it out.
- `webhook_method` is ignored for Telegram (always `POST`).

The delivered message looks like:

```
🔴 Rota pool alert
Pool: US Residential (#1)
Active proxies: 8 / 9
Threshold: 9
```

### Per-User Routing

1. Create pools for each location/use-case
2. Go to **Proxy Users**, click **Add User**
3. Set a main pool and optional fallback pools (in priority order)
4. Configure max retries across the chain and an optional `requests_per_minute` cap

Users connect as:
```
http://username:password@your-proxy-host:8000
```

If the main pool has no live IPs the request automatically cascades to the next fallback pool.

---

## 🔐 API Authentication

All API endpoints require a JWT bearer token obtained from `POST /api/v1/auth/login`.

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"yourpassword"}' | jq -r '.token')

# Use token
curl -H "Authorization: Bearer $TOKEN" http://localhost/api/v1/proxies
```

Public endpoints (no token required):
- `GET /health`
- `POST /api/v1/auth/login`

### Brute-Force Protection

The login endpoint has two independent rate-limit mechanisms:

| Mechanism | Trigger | Response |
|-----------|---------|----------|
| **Per-IP block** | ≥ `AUTH_IP_MAX_ATTEMPTS` failed attempts from one IP within `AUTH_IP_WINDOW_MINUTES` minutes | `429` — IP blocked for `AUTH_IP_BLOCK_MINUTES` minutes |
| **Global lockout** | ≥ `AUTH_GLOBAL_MAX_PER_MINUTE` total attempts per minute across all IPs | `429` — login disabled for everyone for `AUTH_GLOBAL_LOCKOUT_MINUTES` minute(s) |

Both responses include a `Retry-After` header. All thresholds are configurable via `.env`.

> **Behind a reverse proxy?** Per-IP tracking uses the socket peer address by
> default. Set `TRUST_PROXY_HEADERS=true` (the default in the bundled Caddy
> setup) so the real client IP is read from `X-Forwarded-For` instead of the
> proxy's address. Leave it `false` when the API is exposed directly, otherwise
> a client could forge the header to dodge the per-IP block.

The dashboard automatically redirects to the login page with a *"Session expired"* message when a `401` response is received.

---

## 🤝 Contributing

Contributions are welcome! We appreciate meaningful contributions that add value to the project.

### How to Contribute

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Commit your changes**: `git commit -m 'Add amazing feature'`
5. **Push to the branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

### Contribution Guidelines

- Write clear, descriptive commit messages
- Add tests for new features
- Update documentation as needed
- Follow existing code style and conventions
- Ensure all tests pass before submitting PR
- One feature/fix per pull request

**Note**: Pull requests that do not contribute significant improvements or fixes will not be accepted.

### Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/my-feature

# 2. Make changes and test
make test

# 3. Commit changes
git add .
git commit -m "feat: add my feature"

# 4. Push and create PR
git push origin feature/my-feature
```

---

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <p>
    <sub>Built with ❤️ by <a href="https://github.com/alpkeskin">Alp Keskin</a></sub>
  </p>
  <p>
    <sub>⭐ Star this repository if you find it useful!</sub>
  </p>
</div>

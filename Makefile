# Rota — common tasks. Run `make` or `make help` to list targets.

.DEFAULT_GOAL := help

.PHONY: help up build down restart logs ps password dev-core dev-dashboard

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

up: ## Start the full stack (Docker) — open http://localhost
	docker compose up -d

build: ## Rebuild images and start
	docker compose up -d --build

down: ## Stop and remove the stack
	docker compose down

restart: ## Restart all services
	docker compose restart

logs: ## Tail logs from all services
	docker compose logs -f

ps: ## Show service status
	docker compose ps

password: ## Print the first-boot admin password from the core logs
	@docker compose logs rota-core 2>/dev/null | grep -i "Generated admin password" \
		|| echo "No generated password found (you set ROTA_ADMIN_PASSWORD, or it's already seeded)."

dev-core: ## Run the Go core locally (needs a running TimescaleDB)
	cd core && go run ./cmd/server

dev-dashboard: ## Run the Next.js dashboard locally (pnpm dev)
	cd dashboard && pnpm install && pnpm dev

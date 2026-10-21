.PHONY: up down build logs test key help

up: ## Start all containers
	docker compose up -d --build

down: ## Stop and remove containers
	docker compose down

build: ## Rebuild the app image
	docker compose build app

logs: ## Tail the app log
	docker compose logs -f app

test: ## Run vet and the tests
	go vet ./...
	go test ./...

key: ## Print a new SESSION_KEY
	@openssl rand -base64 32

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

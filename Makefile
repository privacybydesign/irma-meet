.PHONY: up down build shell migrate seed key logs

up: ## Start all containers
	docker compose up -d

down: ## Stop and remove containers
	docker compose down

build: ## Rebuild the app image
	docker compose build app

shell: ## Open a shell in the app container
	docker compose exec app bash

key: ## Generate a new APP_KEY and write it into .env
	@KEY=$$(docker compose exec app php -r "echo 'base64:'.base64_encode(random_bytes(32));"); \
	sed -i '' "s|^APP_KEY=.*|APP_KEY=$$KEY|" .env && echo "APP_KEY set."

migrate: ## Run database migrations
	docker compose exec app php artisan migrate

seed: ## Seed the database
	docker compose exec app php artisan db:seed

fresh: ## Drop all tables and re-run migrations + seeds
	docker compose exec app php artisan migrate:fresh --seed

logs: ## Tail Laravel log
	docker compose exec app tail -f storage/logs/laravel.log

composer: ## Run composer (pass CMD="install", "update", etc.)
	docker compose exec app composer $(CMD)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

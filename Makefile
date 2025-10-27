# Helper commands

help: ## Show this help message
	@echo "TPS Benchmark Demo - Available Commands:"
	@echo "----------------------------------"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Development environment commands

dev-up: ## Start development environment in detached mode
	docker compose -f docker-compose.dev.yml up -d

dev-down: ## Stop development environment
	docker compose -f docker-compose.dev.yml down -v

dev-logs: ## Show logs from all services
	docker compose -f docker-compose.dev.yml logs -f

dev-status: ## Show status of all services
	docker compose -f docker-compose.dev.yml ps

dev-restart: ## Restart development environment
	docker compose -f docker-compose.dev.yml restart

# Testing commands

test-python: ## Test Python stack (py-gateway)
	@echo "Testing Python stack..."
	curl -X POST "http://localhost:4001/test/load/burst" \
		-H "Content-Type: application/json" \
		-d '{"service_name": "py-gateway", "protocol": "bl2", "total_reqs": 1, "payload": {"account_number": "1001000000001"}}'

test-go: ## Test Go stack (go-gateway)
	@echo "Testing Go stack..."
	curl -X POST "http://localhost:4001/test/load/burst" \
		-H "Content-Type: application/json" \
		-d '{"service_name": "go-gateway", "protocol": "grpc", "total_reqs": 1, "payload": {"account_number": "1001000000001"}}'

# Development tools

sqlc: ## Generate SQLC code
	cd go-core/store/postgres_store && sqlc generate

build: ## Build all Docker images
	docker compose -f docker-compose.dev.yml build

.PHONY: help dev-up dev-down dev-logs dev-status dev-restart test-python test-go sqlc build
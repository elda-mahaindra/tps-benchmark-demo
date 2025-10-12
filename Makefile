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

.PHONY: dev-up dev-down help
# Helper commands

help: ## Show this help message
	@echo "Available Commands:"
	@echo "--------------------------------------------------------------------"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Development environment commands

go-gateway-pb: ## Generate go-gateway protobuf files
	protoc --proto_path=adapter/go_gateway_adapter/pb adapter/go_gateway_adapter/pb/*.proto --go_out=adapter/go_gateway_adapter/pb --go_opt=paths=source_relative --go-grpc_out=adapter/go_gateway_adapter/pb --go-grpc_opt=paths=source_relative

run-dev: ## Run the development server
	go run cmd/*.go start

PHONY: help run-dev
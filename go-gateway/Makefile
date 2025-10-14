# Helper commands

help: ## Show this help message
	@echo "Available Commands:"
	@echo "--------------------------------------------------------------------"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Development environment commands

go-switching-pb: ## Generate go-switching protobuf files
	protoc --proto_path=adapter/go_switching_adapter/pb adapter/go_switching_adapter/pb/*.proto --go_out=adapter/go_switching_adapter/pb --go_opt=paths=source_relative --go-grpc_out=adapter/go_switching_adapter/pb --go-grpc_opt=paths=source_relative

pb: ## Generate protobuf files
	protoc --proto_path=api/grpc_api/pb api/grpc_api/pb/*.proto --go_out=api/grpc_api/pb --go_opt=paths=source_relative --go-grpc_out=api/grpc_api/pb --go-grpc_opt=paths=source_relative

run-dev: ## Run the development server
	go run cmd/*.go start

PHONY: help go-switching-pb pb run-dev
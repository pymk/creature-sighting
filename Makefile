.PHONY: help run test fmt lint clean

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run: ## Run the server
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

lint: ## Run linters
	golangci-lint run

clean: ## Clean build artifacts
	go clean
	rm -f creature-sighting
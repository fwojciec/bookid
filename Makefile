.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-v
test-v: ## Run tests with verbose output
	go test -v ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix ./...

.PHONY: tools
tools: ## Install development tools
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: validate
validate: ## Run all validation checks (format, vet, test, lint)
	@echo "Running validation checks..."
	@echo "1. Formatting..."
	@go fmt ./...
	@echo "2. Vetting..."
	@go vet ./...
	@echo "3. Testing..."
	@go test ./...
	@echo "4. Tidying modules..."
	@go mod tidy
	@echo "5. Linting..."
	@golangci-lint run ./...
	@echo "✅ All validation checks passed!"

.PHONY: clean
clean: ## Clean build artifacts
	go clean -cache
	rm -f bookid
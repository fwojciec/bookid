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
validate: ## Run all validation checks (same as CI)
	@echo "Running validation checks..."
	@echo "1. Checking formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "❌ Code is not formatted. Run 'go fmt ./...'"; \
		gofmt -d .; \
		exit 1; \
	fi
	@echo "✅ Code formatting OK"
	@echo "2. Running go vet..."
	@go vet ./...
	@echo "✅ go vet OK"
	@echo "3. Running tests with race detector..."
	@go test -race -shuffle=on ./...
	@echo "✅ Tests OK"
	@echo "4. Checking go.mod..."
	@cp go.mod go.mod.bak
	@cp go.sum go.sum.bak 2>/dev/null || true
	@go mod tidy
	@if ! diff -q go.mod go.mod.bak >/dev/null 2>&1 || ! diff -q go.sum go.sum.bak >/dev/null 2>&1; then \
		echo "❌ go.mod/go.sum needs updating. Run 'go mod tidy'"; \
		diff go.mod.bak go.mod || true; \
		diff go.sum.bak go.sum || true; \
		rm -f go.mod.bak go.sum.bak; \
		exit 1; \
	fi
	@rm -f go.mod.bak go.sum.bak
	@echo "✅ go.mod OK"
	@echo "5. Running linter..."
	@golangci-lint run --timeout=5m ./...
	@echo "✅ Linting OK"
	@echo ""
	@echo "✅ All validation checks passed!"

.PHONY: clean
clean: ## Clean build artifacts
	go clean -cache
	rm -f bookid
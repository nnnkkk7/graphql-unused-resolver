.PHONY: build test run clean install help

BINARY_NAME=graphql-unused-resolver
BIN_DIR=bin
CMD_DIR=cmd/$(BINARY_NAME)

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run with test data"
	@echo "  run        - Build and run with test data"
	@echo "  clean      - Remove built binaries"
	@echo "  install    - Install the binary to GOPATH/bin"
	@echo "  fmt        - Format Go code"
	@echo "  vet        - Run go vet"
	@echo "  mod        - Tidy go modules"
	@echo "  test-unit  - Run unit tests with coverage"
	@echo "  lint       - Run golangci-lint"
	@echo "  fmt-check  - Check if code is formatted"
	@echo "  mod-check  - Check if go.mod and go.sum are tidy"
	@echo "  ci         - Run all CI checks locally"
	@echo "  help       - Display this help message"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build complete: $(BIN_DIR)/$(BINARY_NAME)"

## test: Run with test data
test: build
	@echo "Running analysis on test data..."
	@./$(BIN_DIR)/$(BINARY_NAME) \
		--schema testdata/simple/schema.graphql \
		--resolvers testdata/simple/resolvers

## run: Build and run with test data
run: build test

## clean: Remove built binaries
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "✅ Clean complete"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(CMD_DIR)
	@echo "✅ Install complete"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete"

## mod: Tidy go modules
mod:
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "✅ Modules tidied"

## test-unit: Run unit tests with coverage
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "✅ Tests complete"

## lint: Run golangci-lint
lint:
	@echo "Running linters..."
	@golangci-lint run
	@echo "✅ Lint complete"

## fmt-check: Check if code is formatted
fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "❌ Code is not formatted. Run 'make fmt' to fix." && gofmt -l . && exit 1)
	@echo "✅ Code is properly formatted"

## mod-check: Check if go.mod and go.sum are tidy
mod-check:
	@echo "Checking go modules..."
	@cp go.mod go.mod.bak
	@cp go.sum go.sum.bak
	@go mod tidy
	@diff go.mod go.mod.bak > /dev/null || (echo "❌ go.mod is not tidy. Run 'make mod' to fix." && rm go.mod.bak go.sum.bak && exit 1)
	@diff go.sum go.sum.bak > /dev/null || (echo "❌ go.sum is not tidy. Run 'make mod' to fix." && rm go.mod.bak go.sum.bak && exit 1)
	@rm go.mod.bak go.sum.bak
	@echo "✅ Modules are tidy"

## ci: Run all CI checks locally
ci: fmt-check mod-check vet lint test-unit build
	@echo "✅ All CI checks passed"

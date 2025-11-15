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

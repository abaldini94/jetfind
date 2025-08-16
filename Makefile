# Makefile
BINARY_NAME=jetfind
BUILD_DIR=./bin
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)
MAIN_PACKAGE=./cmd/jetfind
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all build run test test-verbose clean lint install help

all: build

build:
	@echo "==> Building... (version $(VERSION))..."
	@go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "==> Exe file created $(BINARY_PATH)"

run: build
	@echo "==> Running from build.."
	@$(BINARY_PATH)

test:
	@echo "==> Running tests..."
	@go test ./...

test-verbose:
	@echo "==> Running tests (verbose)..."
	@go test -v ./...

clean:
	@echo "==> Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@echo "==> Pulizia completata."

lint:
	@echo "==> linter execution..."
	@golangci-lint run

install:
	@echo "==> Installation of $(BINARY_NAME)..."
	@go install -ldflags="-X main.version=$(VERSION)" $(MAIN_PACKAGE)
	@echo "==> $(BINARY_NAME) installed in $(shell go env GOPATH)/bin"

# Show available commands
help:
	@echo "Comandi disponibili:"
	@echo "  make build         - Compile the code"
	@echo "  make run           - Run the code "
	@echo "  make test          - Run the tests"
	@echo "  make test-verbose  - Run the test in verbose mode"
	@echo "  make clean         - Remove compiled files"
	@echo "  make lint          - Linter"
	@echo "  make install       - Install globally with go install"
	@echo "  make help          - Show this help message"

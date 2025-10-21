.PHONY: help build install install-user uninstall clean test run fmt vet

# Binary configuration
BINARY_NAME=translation-key-usage-tracker
MAIN_FILE=main.go

# Installation directories
SYSTEM_INSTALL_DIR=/usr/local/bin
USER_INSTALL_DIR=$(HOME)/bin

# Build output directory
BUILD_DIR=.

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod

# Default target
help:
	@echo "Translation Key Usage Tracker - Makefile commands:"
	@echo ""
	@echo "  make build         - Build the binary"
	@echo "  make install       - Install to $(SYSTEM_INSTALL_DIR) (requires sudo)"
	@echo "  make install-user  - Install to $(USER_INSTALL_DIR) (no sudo needed)"
	@echo "  make uninstall     - Remove from $(SYSTEM_INSTALL_DIR)"
	@echo "  make clean         - Remove built binaries"
	@echo "  make test          - Run tests"
	@echo "  make run           - Build and run with example args"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make deps          - Download dependencies"
	@echo ""

# Check if Go is installed
check-go:
		@which $(GOCMD) > /dev/null || (echo "❌ Error: Go is not installed" && echo "Install Go from: https://go.dev/doc/install" && echo "Or use: brew install go (macOS)" && exit 1)
		@echo "✅ Go is installed: $$($(GOCMD) version)"

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install to system directory (requires sudo)
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(SYSTEM_INSTALL_DIR)..."
	@if [ -w "$(SYSTEM_INSTALL_DIR)" ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(SYSTEM_INSTALL_DIR)/; \
	else \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(SYSTEM_INSTALL_DIR)/; \
	fi
	@echo "✅ Installation complete!"
	@echo ""
	@echo "Usage: $(BINARY_NAME) --help"

# Install to user directory (no sudo needed)
install-user: build
	@echo "📦 Installing $(BINARY_NAME) to $(USER_INSTALL_DIR)..."
	@mkdir -p $(USER_INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(USER_INSTALL_DIR)/
	@echo "✅ Installation complete!"
	@echo ""
	@echo "⚠️  Make sure $(USER_INSTALL_DIR) is in your PATH"
	@echo "Add this to your ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\$$PATH:$(USER_INSTALL_DIR)"
	@echo ""
	@echo "Usage: $(BINARY_NAME) --help"

# Uninstall from system directory
uninstall:
	@echo "🗑️  Removing $(BINARY_NAME)..."
	@if [ -w "$(SYSTEM_INSTALL_DIR)" ]; then \
		rm -f $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo rm -f $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@rm -f $(USER_INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Uninstall complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -f output.json output.csv
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	@$(GOTEST) -v ./...

# Format code
fmt:
	@echo "✨ Formatting code..."
	@$(GOFMT) ./...

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	@$(GOVET) ./...

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "✅ Dependencies ready"

# Run with example arguments (for development)
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@if [ -z "$(XXL_FES_PATH)" ]; then \
		echo "❌ Error: XXL_FES_PATH environment variable not set"; \
		echo "Set it with: export XXL_FES_PATH=/path/to/project"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) --translations translations.yaml --output output.json --format json

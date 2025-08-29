# Air Quality Monitor Makefile

BINARY_NAME=air-quality-monitor
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: all build clean run test server help

all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for multiple platforms
build-all: build-linux build-windows build-darwin

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

# Run the application in command-line mode
run:
	@echo "Running $(BINARY_NAME) in command-line mode..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom device URL
run-custom:
	@echo "Running $(BINARY_NAME) with custom device URL..."
	@./$(BUILD_DIR)/$(BINARY_NAME) $(DEVICE_URL)

# Run the web server
server:
	@echo "Starting $(BINARY_NAME) web server..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --server

# Run server with custom settings
server-custom:
	@echo "Starting $(BINARY_NAME) web server with custom settings..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --server $(DEVICE_URL) $(SERVER_ADDR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install the binary to system path
install:
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Show help
help:
	@echo "Air Quality Monitor - Available targets:"
	@echo ""
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for Linux, Windows, and macOS"
	@echo "  run            - Run in command-line mode"
	@echo "  run-custom     - Run with custom device URL (set DEVICE_URL)"
	@echo "  server         - Start web server"
	@echo "  server-custom  - Start server with custom settings (set DEVICE_URL, SERVER_ADDR)"
	@echo "  deps           - Install dependencies"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install to /usr/local/bin"
	@echo "  uninstall      - Remove from /usr/local/bin"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run-custom DEVICE_URL=http://192.168.1.150/json"
	@echo "  make server-custom DEVICE_URL=http://192.168.1.150/json SERVER_ADDR=:9090"

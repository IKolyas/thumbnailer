# Variables
BINARY_NAME ?= previewer
DOCKER_IMAGE ?= previewer
CONFIG_PATH ?= configs/config.json
SRC_HOST ?= source-host.local:8081

.PHONY: build run lint clean docker-build docker-run docker-stop server-run server-stop unit-test integration-test

# Build the application
build:
	@echo "Building application..."
	go build -o ./bin/$(BINARY_NAME) ./cmd/previewer
	@echo "Build complete"

# Build and run
run: build
	@echo "Starting application..."
	./bin/$(BINARY_NAME) --config=$(CONFIG_PATH)

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean complete"

# Docker operations
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile . --no-cache

docker-run:
	@echo "Starting Docker container..."
	docker compose -f deployments/docker-compose.yaml up previewer -d --build

docker-stop:
	@echo "Stopping Docker container..."
	docker compose -f deployments/docker-compose.yaml down previewer

# Test server operations
server-run:
	docker compose -f deployments/docker-compose.yaml up storage-server -d --build

server-stop:
	docker compose -f deployments/docker-compose.yaml down storage-server

# Tests
unit-test: 
	go test -v ./internal/logger/
	go test -v ./internal/storage/memory

integration-test: server-run docker-run
	go test ./integrations
	$(MAKE) server-stop
	$(MAKE) docker-stop
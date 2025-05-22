# Variables
BINARY_NAME ?= previewer
DOCKER_IMAGE ?= previewer
CONFIG_PATH ?= configs/config.json
SRC_HOST ?= localhost

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
	docker compose -f deployments/docker-compose.yaml up -d --build

docker-stop:
	@echo "Stopping Docker container..."
	docker compose -f deployments/docker-compose.yaml down

# Test server operations
server-run:
	docker compose -f deployments/server/docker-compose.yaml up -d --build

server-stop:
	docker compose -f deployments/server/docker-compose.yaml down

# Tests
unit-test: 
	go test -v ./internal/logger/
	go test -v ./internal/storage/memory

integration-test: server-run docker-run
	SRC_HOST=$(SRC_HOST) go test ./integrations
	$(MAKE) server-stop
	$(MAKE) docker-stop
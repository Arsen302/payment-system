.PHONY: setup build proto run-auth run-payment run-notification run-all docker-up docker-down clean test help

# Default environment file to use
ENV_FILE ?= .env.local

help:
	@echo "Payment Processing System Make Commands"
	@echo "======================================="
	@echo "setup        - Install dependencies for all services"
	@echo "build        - Build all services"
	@echo "proto        - Generate gRPC code from proto files"
	@echo "run-auth     - Run Auth Service"
	@echo "run-payment  - Run Payment Service"
	@echo "run-notification - Run Notification Service"
	@echo "run-all      - Run all services (requires tmux)"
	@echo "docker-up    - Start all services using Docker Compose"
	@echo "docker-down  - Stop Docker Compose services"
	@echo "clean        - Clean build artifacts"
	@echo "test         - Run tests in all services"
	@echo ""
	@echo "Example: make run-auth ENV_FILE=.env.local"

setup:
	@echo "Installing dependencies for all services..."
	cd auth-service && go mod tidy
	cd payment-service && go mod tidy
	cd notification-service && go mod tidy

build:
	@echo "Building all services..."
	cd auth-service && go build -o bin/auth-service ./cmd/server
	cd payment-service && go build -o bin/payment-service ./cmd/server
	cd notification-service && go build -o bin/notification-service ./cmd/server

proto:
	@echo "Generating gRPC code from proto files..."
	mkdir -p auth-service/api/proto/auth
	mkdir -p payment-service/api/proto/payment
	protoc --go_out=. --go-grpc_out=. proto/auth.proto
	protoc --go_out=. --go-grpc_out=. proto/payment.proto

run-auth:
	@echo "Running Auth Service..."
	cd auth-service && ENV_FILE=$(ENV_FILE) go run cmd/server/main.go

run-payment:
	@echo "Running Payment Service..."
	cd payment-service && ENV_FILE=$(ENV_FILE) go run cmd/server/main.go

run-notification:
	@echo "Running Notification Service..."
	cd notification-service && ENV_FILE=$(ENV_FILE) go run cmd/server/main.go

run-all:
	@echo "Running all services (requires tmux)..."
	tmux new-session -d -s payment-system "make run-auth ENV_FILE=$(ENV_FILE)"
	tmux split-window -h -t payment-system "make run-payment ENV_FILE=$(ENV_FILE)"
	tmux split-window -v -t payment-system "make run-notification ENV_FILE=$(ENV_FILE)"
	tmux select-layout -t payment-system tiled
	tmux attach-session -t payment-system

docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

clean:
	@echo "Cleaning build artifacts..."
	rm -rf auth-service/bin
	rm -rf payment-service/bin
	rm -rf notification-service/bin

test:
	@echo "Running tests in all services..."
	cd auth-service && go test -v ./...
	cd payment-service && go test -v ./...
	cd notification-service && go test -v ./... 
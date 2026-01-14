.PHONY: help install dev build clean test lint frontend-install frontend-dev frontend-build frontend-clean

help:
	@echo "Available commands:"
	@echo "  make install          - Install all dependencies (Go + frontend)"
	@echo "  make dev              - Run application in development mode"
	@echo "  make build            - Build production binary"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make test             - Run Go tests"
	@echo "  make lint             - Run linters"
	@echo "  make frontend-install - Install frontend dependencies only"
	@echo "  make frontend-dev     - Run frontend in development mode"
	@echo "  make frontend-build   - Build frontend only"
	@echo "  make frontend-clean   - Clean frontend build artifacts"

install: frontend-install
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

dev:
	@echo "Starting Wails development mode..."
	wails dev

build:
	@echo "Building production binary..."
	wails build

clean:
	@echo "Cleaning build artifacts..."
	rm -rf build/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite

test:
	@echo "Running Go tests..."
	go test -v ./...

lint:
	@echo "Running Go linters..."
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi
	@echo "Running frontend linters..."
	cd frontend && npm run lint || echo "No lint script configured"

frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

frontend-dev:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

frontend-build:
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-clean:
	@echo "Cleaning frontend artifacts..."
	rm -rf frontend/dist
	rm -rf frontend/node_modules

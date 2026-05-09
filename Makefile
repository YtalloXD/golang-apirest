.PHONY: help build run test clean fmt lint

help:
	@echo "Video Games REST API - Makefile Commands"
	@echo "=========================================="
	@echo "build      - Build the application"
	@echo "run        - Run the application"
	@echo "test       - Run tests"
	@echo "fmt        - Format code"
	@echo "lint       - Run linter"
	@echo "clean      - Remove build artifacts"
	@echo "deps       - Download dependencies"
	@echo "help       - Display this help message"

build:
	go build -o video-games-api.exe -v

run:
	go run main.go

test:
	go test -v ./...

fmt:
	go fmt ./...
	goimports -w .

lint:
	golint ./...
	go vet ./...

clean:
	rm -f video-games-api.exe
	go clean

deps:
	go mod download
	go mod tidy

# Notes App Backend Makefile
# Detect operating system
ifeq ($(OS),Windows_NT)
    BINARY_NAME = bin/main.exe
    RM = rmdir /s /q
else
    BINARY_NAME = bin/main
    RM = rm -rf
endif

# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Map environment variables from .env to application-expected variables
export DB_USERNAME ?= $(DB_USER)

.PHONY: all help run build clean db-up db-down db-logs tidy fmt vet

all: help

help:
	@echo Notes App Backend - Development Makefile
	@echo =========================================
	@echo Usage: make [target]
	@echo.
	@echo Available Targets:
	@echo   help       - Show this help menu with all available commands
	@echo   run        - Start the database container and run the Go application
	@echo   build      - Build the Go application binary for the current OS
	@echo   clean      - Remove build artifacts and binary files
	@echo   db-up      - Start the PostgreSQL database container in detached mode
	@echo   db-down    - Stop and remove the PostgreSQL database container
	@echo   db-logs    - Stream logs from the PostgreSQL database container
	@echo   tidy       - Run go mod tidy to clean up dependencies
	@echo   fmt        - Format all Go source code files
	@echo   vet        - Examine Go source code for suspicious constructs

run:
	@echo Starting the Go application...
	go run cmd/main.go

build: tidy
	@echo Building Go binary to $(BINARY_NAME)...
	@ifeq ($(OS),Windows_NT)
		@if not exist bin mkdir bin
	@else
		@mkdir -p bin
	@endif
	go build -o $(BINARY_NAME) ./cmd/main.go

clean:
	@echo Cleaning build artifacts...
	@ifeq ($(OS),Windows_NT)
		@if exist bin $(RM) bin
	@else
		$(RM) bin
	@endif

db-up:
	@echo Starting PostgreSQL container...
	docker compose up -d

db-down:
	@echo Stopping PostgreSQL container...
	docker compose down

db-logs:
	@echo Streaming database logs...
	docker compose logs -f postgres

tidy:
	@echo Tidying up Go modules...
	go mod tidy

fmt:
	@echo Formatting Go source files...
	go fmt ./...

vet:
	@echo Vetting Go source files...
	go vet ./...

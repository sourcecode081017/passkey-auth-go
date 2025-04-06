# Makefile for Docker Compose application

.PHONY: up down build logs clean restart help

# Default target when just `make` is run
all: up

# Start the application using docker-compose with build flag
up:
	@echo "Starting application containers..."
	docker-compose up --build -d
	@echo "Application is running!"
	@echo "Backend available at: http://localhost:8080"
	@echo "Frontend available at: http://localhost:5173"

# Start the application in foreground mode with logs
start:
	@echo "Starting application containers in foreground mode..."
	docker-compose up --build

# Stop the application
down:
	@echo "Stopping application containers..."
	docker-compose down
	@echo "Application stopped"

# Build the containers without starting
build:
	@echo "Building application containers..."
	docker-compose build

# View logs
logs:
	@echo "Showing logs..."
	docker-compose logs -f

# Clean up all containers, volumes and images
clean:
	@echo "Cleaning up all containers, volumes and images..."
	docker-compose down -v
	@echo "Cleanup complete"

# Restart the application
restart: down up
	@echo "Application restarted"

# Display help information
help:
	@echo "Available targets:"
	@echo "  make              - Start application in detached mode"
	@echo "  make start        - Start application in foreground mode with logs"
	@echo "  make up           - Start application in detached mode"
	@echo "  make down         - Stop application"
	@echo "  make build        - Build containers without starting"
	@echo "  make logs         - View application logs"
	@echo "  make clean        - Remove all containers, volumes and images"
	@echo "  make restart      - Restart application"
	@echo "  make help         - Display this help"
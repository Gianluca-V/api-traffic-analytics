.PHONY: build run test clean setup

# Variables
COMPOSE_FILE = deployments/docker/docker-compose.yml

# Construir todos los microservicios
build:
	@echo "Building microservices..."
	docker-compose -f $(COMPOSE_FILE) build

# Levantar todos los servicios
run:
	@echo "Starting all services..."
	docker-compose -f $(COMPOSE_FILE) up -d

# Detener todos los servicios
stop:
	@echo "Stopping all services..."
	docker-compose -f $(COMPOSE_FILE) down

# Ver logs
logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# Ejecutar tests
test:
	@echo "Running tests..."
	go test ./... -v

# Setup inicial
setup:
	@echo "Setting up project..."
	go mod tidy

# Limpiar
clean:
	@echo "Cleaning up..."
	docker-compose -f $(COMPOSE_FILE) down -v
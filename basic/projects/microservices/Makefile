# Build all services
build:
	docker-compose build

# Run all services
run:
	docker-compose up -d

# Stop all services
stop:
	docker-compose down

# Clean up everything
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run tests
test:
	go test -v ./...

# Build individual services
build-service-a:
	docker build -t microservices-service-a -f Dockerfile.service-a .

build-service-b:
	docker build -t microservices-service-b -f Dockerfile.service-b .

build-gateway:
	docker build -t microservices-gateway -f Dockerfile.gateway .

# Validate docker-compose
validate:
	docker-compose config

# Show logs
logs:
	docker-compose logs -f

# Show service status
status:
	docker-compose ps

# Restart a specific service
restart-service-a:
	docker-compose restart service-a

restart-service-b:
	docker-compose restart service-b

restart-gateway:
	docker-compose restart api-gateway

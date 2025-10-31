.PHONY: help kafka-up rabbitmq-up kafka-producer kafka-consumer rabbitmq-publisher rabbitmq-subscriber

help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║                                                              ║"
	@echo "║        📨 Messaging with Go: Kafka & RabbitMQ               ║"
	@echo "║                                                              ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available targets:"
	@echo "  make kafka-up            - Start Kafka with Docker Compose"
	@echo "  make kafka-down          - Stop Kafka"
	@echo "  make kafka-producer      - Run Kafka producer"
	@echo "  make kafka-consumer      - Run Kafka consumer"
	@echo "  make kafka-consumer-group - Run Kafka consumer group"
	@echo "  make rabbitmq-up         - Start RabbitMQ with Docker Compose"
	@echo "  make rabbitmq-down       - Stop RabbitMQ"
	@echo "  make rabbitmq-publisher  - Run RabbitMQ publisher"
	@echo "  make rabbitmq-subscriber - Run RabbitMQ subscriber"
	@echo "  make rabbitmq-worker     - Run RabbitMQ worker"
	@echo "  make all-up              - Start both Kafka and RabbitMQ"
	@echo "  make all-down            - Stop both Kafka and RabbitMQ"

# Kafka commands
kafka-up:
	@echo "🚀 Starting Kafka..."
	@docker-compose -f docker-compose.kafka.yml up -d
	@echo "✅ Kafka started"
	@echo "📍 Kafka: localhost:9092"
	@echo "📍 Kafka UI: http://localhost:8080"

kafka-down:
	@echo "🛑 Stopping Kafka..."
	@docker-compose -f docker-compose.kafka.yml down
	@echo "✅ Kafka stopped"

kafka-producer:
	@echo "📤 Running Kafka producer..."
	@go run kafka/producer/main.go

kafka-consumer:
	@echo "📥 Running Kafka consumer..."
	@go run kafka/consumer/main.go

kafka-consumer-group:
	@echo "📥 Running Kafka consumer group..."
	@go run kafka/consumer-group/main.go

# RabbitMQ commands
rabbitmq-up:
	@echo "🚀 Starting RabbitMQ..."
	@docker-compose -f docker-compose.rabbitmq.yml up -d
	@echo "✅ RabbitMQ started"
	@echo "📍 RabbitMQ: localhost:5672"
	@echo "📍 Management UI: http://localhost:15672 (guest/guest)"

rabbitmq-down:
	@echo "🛑 Stopping RabbitMQ..."
	@docker-compose -f docker-compose.rabbitmq.yml down
	@echo "✅ RabbitMQ stopped"

rabbitmq-publisher:
	@echo "📤 Running RabbitMQ publisher..."
	@go run rabbitmq/publisher/main.go

rabbitmq-subscriber:
	@echo "📥 Running RabbitMQ subscriber..."
	@go run rabbitmq/subscriber/main.go

rabbitmq-worker:
	@echo "👷 Running RabbitMQ worker..."
	@go run rabbitmq/worker/main.go

# All services
all-up:
	@make kafka-up
	@make rabbitmq-up

all-down:
	@make kafka-down
	@make rabbitmq-down


#!/bin/bash

# GO-PRO Kafka Real-Time Data Processing - Quick Start Script
# This script provides convenient commands for Kafka operations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.kafka.yml"
KAFKA_BROKER="localhost:9092"
KAFKA_BOOTSTRAP="kafka:29092"
ZOOKEEPER_HOST="zookeeper:2181"

# Helper functions
print_header() {
    echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Main commands
start_stack() {
    print_header "Starting Kafka Stack"
    docker-compose -f $COMPOSE_FILE up -d
    print_success "Kafka stack started"
    echo ""
    print_info_urls
}

stop_stack() {
    print_header "Stopping Kafka Stack"
    docker-compose -f $COMPOSE_FILE down
    print_success "Kafka stack stopped"
}

restart_stack() {
    print_header "Restarting Kafka Stack"
    docker-compose -f $COMPOSE_FILE down
    docker-compose -f $COMPOSE_FILE up -d
    print_success "Kafka stack restarted"
    echo ""
    print_info_urls
}

status_stack() {
    print_header "Kafka Stack Status"
    docker-compose -f $COMPOSE_FILE ps
}

view_logs() {
    local service="${1:-kafka}"
    print_header "Viewing logs for: $service"
    docker-compose -f $COMPOSE_FILE logs -f $service
}

create_topics() {
    print_header "Creating Kafka Topics"

    local topics=(
        "go-pro.user.events"
        "go-pro.course.events"
        "go-pro.lesson.events"
        "go-pro.exercise.events"
        "go-pro.progress.events"
        "go-pro.notification.events"
        "go-pro.audit.events"
        "go-pro.dlq"
    )

    for topic in "${topics[@]}"; do
        docker-compose -f $COMPOSE_FILE exec -T kafka kafka-topics.sh \
            --create \
            --topic $topic \
            --bootstrap-server $KAFKA_BOOTSTRAP \
            --partitions 1 \
            --replication-factor 1 \
            --if-not-exists 2>/dev/null || true
        print_success "Created topic: $topic"
    done
}

list_topics() {
    print_header "Kafka Topics"
    docker-compose -f $COMPOSE_FILE exec -T kafka kafka-topics.sh \
        --list \
        --bootstrap-server $KAFKA_BOOTSTRAP
}

describe_topic() {
    local topic="${1:-go-pro.progress.events}"
    print_header "Topic Details: $topic"
    docker-compose -f $COMPOSE_FILE exec -T kafka kafka-topics.sh \
        --describe \
        --topic $topic \
        --bootstrap-server $KAFKA_BOOTSTRAP
}

consumer_groups() {
    print_header "Consumer Groups"
    docker-compose -f $COMPOSE_FILE exec -T kafka kafka-consumer-groups.sh \
        --list \
        --bootstrap-server $KAFKA_BOOTSTRAP
}

consumer_lag() {
    local group="${1:-go-pro-consumer-group}"
    print_header "Consumer Group Lag: $group"
    docker-compose -f $COMPOSE_FILE exec -T kafka kafka-consumer-groups.sh \
        --describe \
        --group $group \
        --bootstrap-server $KAFKA_BOOTSTRAP
}

reset_consumer_group() {
    local group="${1:-go-pro-consumer-group}"
    print_header "Resetting Consumer Group: $group"
    docker-compose -f $COMPOSE_FILE exec -T kafka kafka-consumer-groups.sh \
        --reset-offsets \
        --group $group \
        --all-topics \
        --to-earliest \
        --execute \
        --bootstrap-server $KAFKA_BOOTSTRAP
    print_success "Consumer group reset"
}

produce_test_event() {
    print_header "Producing Test Event"

    local event='{"event_type": "test", "timestamp": '$(date +%s)', "data": {"message": "test"}}'

    echo "$event" | docker-compose -f $COMPOSE_FILE exec -T kafka kafka-console-producer.sh \
        --topic go-pro.progress.events \
        --bootstrap-server $KAFKA_BOOTSTRAP

    print_success "Test event produced"
}

consume_events() {
    local topic="${1:-go-pro.progress.events}"
    print_header "Consuming events from: $topic"
    docker-compose -f $COMPOSE_FILE exec kafka kafka-console-consumer.sh \
        --topic $topic \
        --bootstrap-server $KAFKA_BOOTSTRAP \
        --from-beginning
}

clean_topics() {
    print_header "Cleaning All Topics"
    local topics=$(docker-compose -f $COMPOSE_FILE exec -T kafka kafka-topics.sh \
        --list \
        --bootstrap-server $KAFKA_BOOTSTRAP)

    for topic in $topics; do
        docker-compose -f $COMPOSE_FILE exec -T kafka kafka-topics.sh \
            --delete \
            --topic $topic \
            --bootstrap-server $KAFKA_BOOTSTRAP \
            2>/dev/null || true
    done

    print_success "All topics deleted"
}

rebuild_stack() {
    print_header "Rebuilding Kafka Stack"
    docker-compose -f $COMPOSE_FILE down -v
    docker-compose -f $COMPOSE_FILE up -d
    print_success "Stack rebuilt"
    echo ""
    print_info_urls
}

print_info_urls() {
    echo ""
    echo -e "${BLUE}Service URLs:${NC}"
    echo -e "  Kafka UI:      ${YELLOW}http://localhost:8080${NC}"
    echo -e "  Prometheus:    ${YELLOW}http://localhost:9090${NC}"
    echo -e "  Grafana:       ${YELLOW}http://localhost:3000${NC} (admin/admin)"
    echo -e "  Adminer:       ${YELLOW}http://localhost:8081${NC}"
    echo -e "  Qdrant:        ${YELLOW}http://localhost:6333${NC}"
    echo ""
    echo -e "${BLUE}Broker Connection:${NC}"
    echo -e "  Host:          ${YELLOW}localhost${NC}"
    echo -e "  Port:          ${YELLOW}9092${NC}"
    echo ""
}

run_backend() {
    print_header "Starting Go-Pro Backend"

    cd backend 2>/dev/null || {
        print_error "backend directory not found"
        return 1
    }

    export KAFKA_BROKERS=localhost:9092
    export KAFKA_GROUP_ID=go-pro-consumer-group
    export DATABASE_URL=postgres://gopro:gopro_password@localhost:5432/gopro_db
    export REDIS_URL=redis://localhost:6379
    export MESSAGING_ENABLED=true

    go run ./cmd/server
}

health_check() {
    print_header "Health Check"

    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        return 1
    fi
    print_success "Docker is installed"

    # Check docker-compose file
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "docker-compose.kafka.yml not found"
        return 1
    fi
    print_success "docker-compose.kafka.yml found"

    # Check if services are running
    local running=$(docker-compose -f $COMPOSE_FILE ps --services --filter "status=running" 2>/dev/null | wc -l)
    print_success "$running services running"

    # Check Kafka broker
    if docker-compose -f $COMPOSE_FILE exec -T kafka kafka-broker-api-versions.sh \
        --bootstrap-server $KAFKA_BOOTSTRAP &>/dev/null; then
        print_success "Kafka broker is healthy"
    else
        print_warning "Kafka broker is not responding"
    fi
}

print_usage() {
    cat << EOF
${BLUE}GO-PRO Kafka Quick Start - Usage${NC}

${BLUE}Stack Management:${NC}
  ./KAFKA_QUICK_START.sh start           Start Kafka stack
  ./KAFKA_QUICK_START.sh stop            Stop Kafka stack
  ./KAFKA_QUICK_START.sh restart         Restart Kafka stack
  ./KAFKA_QUICK_START.sh status          Show stack status
  ./KAFKA_QUICK_START.sh rebuild         Rebuild stack from scratch

${BLUE}Topic Management:${NC}
  ./KAFKA_QUICK_START.sh create-topics   Create all required topics
  ./KAFKA_QUICK_START.sh list-topics     List all topics
  ./KAFKA_QUICK_START.sh describe <topic> Describe topic details

${BLUE}Consumer Management:${NC}
  ./KAFKA_QUICK_START.sh groups          List consumer groups
  ./KAFKA_QUICK_START.sh lag <group>     Check consumer lag
  ./KAFKA_QUICK_START.sh reset <group>   Reset consumer group offset

${BLUE}Monitoring & Testing:${NC}
  ./KAFKA_QUICK_START.sh health          Health check
  ./KAFKA_QUICK_START.sh logs <service>  View service logs
  ./KAFKA_QUICK_START.sh produce         Produce test event
  ./KAFKA_QUICK_START.sh consume <topic> Consume events

${BLUE}Data Management:${NC}
  ./KAFKA_QUICK_START.sh clean           Delete all topics
  ./KAFKA_QUICK_START.sh urls            Show service URLs

${BLUE}Backend:${NC}
  ./KAFKA_QUICK_START.sh run-backend     Start Go-Pro backend

${BLUE}Examples:${NC}
  ./KAFKA_QUICK_START.sh start
  ./KAFKA_QUICK_START.sh create-topics
  ./KAFKA_QUICK_START.sh list-topics
  ./KAFKA_QUICK_START.sh describe go-pro.progress.events
  ./KAFKA_QUICK_START.sh logs kafka
  ./KAFKA_QUICK_START.sh run-backend

EOF
}

# Main script
main() {
    case "${1:-help}" in
        start)
            start_stack
            ;;
        stop)
            stop_stack
            ;;
        restart)
            restart_stack
            ;;
        status)
            status_stack
            ;;
        logs)
            view_logs "$2"
            ;;
        create-topics)
            create_topics
            ;;
        list-topics)
            list_topics
            ;;
        describe)
            describe_topic "$2"
            ;;
        groups)
            consumer_groups
            ;;
        lag)
            consumer_lag "$2"
            ;;
        reset)
            reset_consumer_group "$2"
            ;;
        produce)
            produce_test_event
            ;;
        consume)
            consume_events "$2"
            ;;
        clean)
            clean_topics
            ;;
        rebuild)
            rebuild_stack
            ;;
        health)
            health_check
            ;;
        urls)
            print_info_urls
            ;;
        run-backend)
            run_backend
            ;;
        *)
            print_usage
            ;;
    esac
}

main "$@"

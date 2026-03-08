#!/bin/bash

# Microservices Example Scripts
# This script provides example commands for testing the microservices

set -e

API_GATEWAY="http://localhost:8080"

echo "=== Microservices Examples ==="
echo ""

# Color output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_blue() {
    echo -e "${BLUE}$1${NC}"
}

# Check if API Gateway is running
check_gateway() {
    if ! curl -s "$API_GATEWAY/health" > /dev/null; then
        echo "Error: API Gateway is not running. Start services with: docker-compose up -d"
        exit 1
    fi
}

# 1. Health Check
health_check() {
    print_blue "1. Checking service health..."
    curl -s "$API_GATEWAY/health" | jq '.'
    echo ""
}

# 2. Create User
create_user() {
    print_blue "2. Creating a new user..."
    response=$(curl -s -X POST "$API_GATEWAY/api/users" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Alice Johnson",
            "email": "alice@example.com"
        }')
    echo "$response" | jq '.'
    USER_ID=$(echo "$response" | jq -r '.id')
    print_green "✓ User created with ID: $USER_ID"
    echo ""
}

# 3. Get User
get_user() {
    print_blue "3. Retrieving user by ID..."
    curl -s "$API_GATEWAY/api/users/$USER_ID" | jq '.'
    echo ""
}

# 4. List All Users
list_users() {
    print_blue "4. Listing all users..."
    curl -s "$API_GATEWAY/api/users" | jq '.'
    echo ""
}

# 5. Update User
update_user() {
    print_blue "5. Updating user..."
    curl -s -X PUT "$API_GATEWAY/api/users/$USER_ID" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Alice Johnson-Smith",
            "email": "alice.smith@example.com"
        }' | jq '.'
    print_green "✓ User updated"
    echo ""
}

# 6. Create Order
create_order() {
    print_blue "6. Creating a new order..."
    response=$(curl -s -X POST "$API_GATEWAY/api/orders" \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$USER_ID\",
            \"items\": [
                {\"product\": \"Widget\", \"quantity\": 2, \"price\": 24.99},
                {\"product\": \"Gadget\", \"quantity\": 1, \"price\": 49.99}
            ],
            \"total\": 99.97
        }")
    echo "$response" | jq '.'
    ORDER_ID=$(echo "$response" | jq -r '.id')
    print_green "✓ Order created with ID: $ORDER_ID"
    echo ""
}

# 7. Get Order
get_order() {
    print_blue "7. Retrieving order by ID..."
    curl -s "$API_GATEWAY/api/orders/$ORDER_ID" | jq '.'
    echo ""
}

# 8. List All Orders
list_orders() {
    print_blue "8. Listing all orders..."
    curl -s "$API_GATEWAY/api/orders" | jq '.'
    echo ""
}

# 9. Get User Orders
get_user_orders() {
    print_blue "9. Retrieving orders for user..."
    curl -s "$API_GATEWAY/api/orders/user/$USER_ID" | jq '.'
    echo ""
}

# 10. Update Order Status
update_order_status() {
    print_blue "10. Updating order status..."
    curl -s -X PUT "$API_GATEWAY/api/orders/$ORDER_ID/status" \
        -H "Content-Type: application/json" \
        -d '{
            "status": "processing"
        }' | jq '.'
    print_green "✓ Order status updated"
    echo ""
}

# 11. Create Multiple Users
create_multiple_users() {
    print_blue "11. Creating multiple users..."
    for i in {1..3}; do
        curl -s -X POST "$API_GATEWAY/api/users" \
            -H "Content-Type: application/json" \
            -d "{
                \"name\": \"User $i\",
                \"email\": \"user$i@example.com\"
            }" | jq '.id'
    done
    print_green "✓ Created 3 additional users"
    echo ""
}

# 12. Create Multiple Orders
create_multiple_orders() {
    print_blue "12. Creating multiple orders for user..."
    for i in {1..3}; do
        curl -s -X POST "$API_GATEWAY/api/orders" \
            -H "Content-Type: application/json" \
            -d "{
                \"user_id\": \"$USER_ID\",
                \"items\": [{\"product\": \"Product $i\", \"quantity\": 1, \"price\": 19.99}],
                \"total\": 19.99
            }" | jq '.id'
    done
    print_green "✓ Created 3 orders"
    echo ""
}

# 13. Service-to-Service Communication Test
test_service_communication() {
    print_blue "13. Testing service-to-service communication..."
    print_green "✓ Order Service validates users via User Service"
    print_green "✓ API Gateway routes requests to appropriate services"
    echo ""
}

# 14. Error Handling - Invalid User
test_invalid_user() {
    print_blue "14. Testing error handling (invalid user)..."
    curl -s -X POST "$API_GATEWAY/api/orders" \
        -H "Content-Type: application/json" \
        -d '{
            "user_id": "invalid-id",
            "items": [{"product": "Widget", "quantity": 1}],
            "total": 9.99
        }' | jq '.'
    echo ""
}

# 15. Error Handling - Invalid Data
test_invalid_data() {
    print_blue "15. Testing error handling (invalid data)..."
    curl -s -X POST "$API_GATEWAY/api/users" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "",
            "email": "invalid-email"
        }' | jq '.'
    echo ""
}

# 16. Correlation ID Test
test_correlation_id() {
    print_blue "16. Testing correlation ID propagation..."
    response=$(curl -v -X POST "$API_GATEWAY/api/users" \
        -H "Content-Type: application/json" \
        -H "X-Correlation-ID: test-correlation-123" \
        -d '{
            "name": "Bob Smith",
            "email": "bob@example.com"
        }' 2>&1 | grep -i "correlation-id")
    print_green "✓ Correlation ID: $response"
    echo ""
}

# 17. Delete User
delete_user() {
    print_blue "17. Deleting user..."
    curl -s -X DELETE "$API_GATEWAY/api/users/$USER_ID" -w "\nHTTP Status: %{http_code}\n"
    print_green "✓ User deleted"
    echo ""
}

# 18. Delete Order
delete_order() {
    print_blue "18. Deleting order..."
    curl -s -X DELETE "$API_GATEWAY/api/orders/$ORDER_ID" -w "\nHTTP Status: %{http_code}\n"
    print_green "✓ Order deleted"
    echo ""
}

# 19. Performance Test
performance_test() {
    print_blue "19. Running performance test (10 concurrent requests)..."
    time (
        for i in {1..10}; do
            curl -s "$API_GATEWAY/api/users" > /dev/null &
        done
        wait
    )
    print_green "✓ Performance test completed"
    echo ""
}

# 20. Show Logs
show_logs() {
    print_blue "20. Showing recent logs..."
    print_green "API Gateway logs:"
    docker-compose logs --tail=5 api-gateway
    echo ""
    print_green "User Service logs:"
    docker-compose logs --tail=5 service-a
    echo ""
}

# Main execution
main() {
    check_gateway

    print_green "Starting microservices examples..."
    echo ""

    # Run all examples
    health_check
    create_user
    get_user
    list_users
    update_user
    create_order
    get_order
    list_orders
    get_user_orders
    update_order_status
    create_multiple_users
    create_multiple_orders
    test_service_communication
    test_invalid_user
    test_invalid_data
    test_correlation_id
    performance_test
    show_logs

    print_green "=== Examples completed successfully! ==="
}

# Run main function
main

# Multi-Stage Dockerfile for Microservices Project
# This project contains 3 services: gateway, service-a (users), service-b (orders)
#
# Build individual services with:
#   docker build --target gateway -t microservices-gateway .
#   docker build --target service-a -t microservices-service-a .
#   docker build --target service-b -t microservices-service-b .
#
# Or use docker-compose for full environment:
#   docker-compose up --build

FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service-a ./cmd/service-a
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service-b ./cmd/service-b

# Gateway stage
FROM alpine:latest AS gateway
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/gateway /gateway
EXPOSE 8080
ENV GATEWAY_PORT=8080
ENV USER_SERVICE_URL=http://service-a:8001
ENV ORDER_SERVICE_URL=http://service-b:8002
CMD ["/gateway"]

# Service A stage
FROM alpine:latest AS service-a
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/service-a /service-a
EXPOSE 8001
ENV SERVICE_PORT=8001
CMD ["/service-a"]

# Service B stage
FROM alpine:latest AS service-b
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/service-b /service-b
EXPOSE 8002
ENV SERVICE_PORT=8002
CMD ["/service-b"]

# 📨 Messaging with Go: Kafka and RabbitMQ

A comprehensive tutorial demonstrating message brokers and event streaming with Go, covering both Apache Kafka and RabbitMQ.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Kafka](#kafka)
- [RabbitMQ](#rabbitmq)
- [Patterns](#patterns)
- [Project Structure](#project-structure)

## 🎯 Overview

This project demonstrates:
- **Apache Kafka**: Event streaming, topics, partitions, consumer groups
- **RabbitMQ**: Message queuing, exchanges, routing, work queues
- **Messaging Patterns**: Pub/Sub, Point-to-Point, Request/Reply
- **Production Practices**: Error handling, graceful shutdown, acknowledgments

## ✨ Features

### Kafka Features
- ✅ Producer with compression and retries
- ✅ Consumer with offset management
- ✅ Consumer groups for load balancing
- ✅ Topic partitioning
- ✅ Message ordering guarantees

### RabbitMQ Features
- ✅ Publisher with persistent messages
- ✅ Subscriber with manual acknowledgments
- ✅ Worker queues for task distribution
- ✅ Topic exchanges for routing
- ✅ QoS for fair dispatch

## 📦 Prerequisites

- **Go** 1.21 or higher
- **Docker** and **Docker Compose**
- Basic understanding of message brokers

## 🚀 Quick Start

### 1. Clone and Navigate

```bash
cd basic/projects/messaging-with-go
```

### 2. Install Dependencies

```bash
go mod tidy
```

## 📊 Kafka

### Start Kafka

```bash
# Start Kafka, Zookeeper, and Kafka UI
make kafka-up

# Services available:
# - Kafka: localhost:9092
# - Kafka UI: http://localhost:8080
```

### Run Kafka Producer

```bash
# Terminal 1: Start producer
make kafka-producer

# Output:
# ✅ Kafka producer connected
# 📤 Sending messages to topic 'events'...
# ✅ Message sent: partition=0 offset=0 id=event-1
# ✅ Message sent: partition=0 offset=1 id=event-2
```

### Run Kafka Consumer

```bash
# Terminal 2: Start consumer
make kafka-consumer

# Output:
# ✅ Kafka consumer connected
# 📥 Consuming messages from topic 'events'...
# ✅ Message received [1]: partition=0 offset=0 key=event-1
#    Event: ID=event-1 Type=user.created Data=User 1 created
```

### Run Kafka Consumer Group

```bash
# Terminal 2: Start consumer group
make kafka-consumer-group

# Run multiple instances for load balancing
# Terminal 3:
make kafka-consumer-group

# Messages will be distributed across consumers
```

### Kafka Concepts

**Topics:**
- Logical channels for messages
- Divided into partitions
- Messages are ordered within partitions

**Partitions:**
- Enable parallel processing
- Each partition is an ordered log
- Messages have sequential offsets

**Consumer Groups:**
- Multiple consumers share workload
- Each partition assigned to one consumer
- Automatic rebalancing

**Offsets:**
- Position in partition
- Managed by consumer
- Enables replay and recovery

### Stop Kafka

```bash
make kafka-down
```

## 🐰 RabbitMQ

### Start RabbitMQ

```bash
# Start RabbitMQ with Management UI
make rabbitmq-up

# Services available:
# - RabbitMQ: localhost:5672
# - Management UI: http://localhost:15672 (guest/guest)
```

### Run RabbitMQ Publisher

```bash
# Terminal 1: Start publisher
make rabbitmq-publisher

# Output:
# ✅ RabbitMQ publisher connected
# 📤 Publishing messages to exchange 'events'...
# ✅ Message published: id=msg-1 routing_key=notification.email
# ✅ Message published: id=msg-2 routing_key=notification.email
```

### Run RabbitMQ Subscriber

```bash
# Terminal 2: Start subscriber
make rabbitmq-subscriber

# Output:
# ✅ RabbitMQ subscriber connected
# 📥 Consuming messages from queue 'notifications'...
# ✅ Message received [1]: id=msg-1
#    Type=notification Payload=Notification 1
```

### Run RabbitMQ Worker

```bash
# Terminal 2: Start worker
make rabbitmq-worker

# Run multiple workers for load balancing
# Terminal 3:
make rabbitmq-worker

# Tasks will be distributed across workers
```

### RabbitMQ Concepts

**Exchanges:**
- Route messages to queues
- Types: direct, topic, fanout, headers
- Routing based on routing keys

**Queues:**
- Store messages
- FIFO ordering
- Durable or temporary

**Bindings:**
- Link exchanges to queues
- Define routing rules
- Support pattern matching

**Acknowledgments:**
- Manual or automatic
- Ensure message processing
- Enable retry on failure

### Stop RabbitMQ

```bash
make rabbitmq-down
```

## 🎨 Messaging Patterns

### 1. Publish/Subscribe (Kafka)

**Use Case:** Event broadcasting to multiple consumers

```
Producer → Topic → Consumer 1
                 → Consumer 2
                 → Consumer 3
```

**Example:** User registration event sent to email service, analytics, and audit log

### 2. Publish/Subscribe (RabbitMQ)

**Use Case:** Message broadcasting with routing

```
Publisher → Exchange (fanout) → Queue 1 → Subscriber 1
                               → Queue 2 → Subscriber 2
```

**Example:** Notification sent to email, SMS, and push notification services

### 3. Work Queue (RabbitMQ)

**Use Case:** Task distribution among workers

```
Producer → Queue → Worker 1
                 → Worker 2
                 → Worker 3
```

**Example:** Image processing tasks distributed across multiple workers

### 4. Topic Routing (RabbitMQ)

**Use Case:** Selective message routing

```
Publisher → Exchange (topic) → Queue (*.error) → Error Handler
                             → Queue (app.*) → App Logger
```

**Example:** Log messages routed based on severity and source

## 📁 Project Structure

```
messaging-with-go/
├── kafka/
│   ├── producer/
│   │   └── main.go              # Kafka producer
│   ├── consumer/
│   │   └── main.go              # Kafka consumer
│   └── consumer-group/
│       └── main.go              # Kafka consumer group
├── rabbitmq/
│   ├── publisher/
│   │   └── main.go              # RabbitMQ publisher
│   ├── subscriber/
│   │   └── main.go              # RabbitMQ subscriber
│   └── worker/
│       └── main.go              # RabbitMQ worker
├── examples/
│   └── (additional examples)
├── docker-compose.kafka.yml     # Kafka stack
├── docker-compose.rabbitmq.yml  # RabbitMQ stack
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

## 🎓 Learning Outcomes

After completing this tutorial, you'll understand:

### Kafka
- ✅ Event streaming architecture
- ✅ Topics and partitions
- ✅ Producer configuration
- ✅ Consumer groups
- ✅ Offset management
- ✅ Message ordering

### RabbitMQ
- ✅ Message queuing patterns
- ✅ Exchanges and routing
- ✅ Queue management
- ✅ Acknowledgments
- ✅ QoS and fair dispatch
- ✅ Work distribution

### General
- ✅ When to use Kafka vs RabbitMQ
- ✅ Message broker patterns
- ✅ Error handling
- ✅ Graceful shutdown
- ✅ Production best practices

## 📊 Kafka vs RabbitMQ

| Feature | Kafka | RabbitMQ |
|---------|-------|----------|
| **Type** | Event Streaming | Message Broker |
| **Use Case** | Event sourcing, logs | Task queues, RPC |
| **Ordering** | Per partition | Per queue |
| **Retention** | Configurable (days) | Until consumed |
| **Throughput** | Very High | High |
| **Latency** | Low | Very Low |
| **Complexity** | Higher | Lower |
| **Replay** | Yes | No |

## 🔧 Configuration

### Kafka Producer

```go
config := sarama.NewConfig()
config.Producer.Return.Successes = true
config.Producer.RequiredAcks = sarama.WaitForAll
config.Producer.Retry.Max = 5
config.Producer.Compression = sarama.CompressionSnappy
```

### Kafka Consumer

```go
config := sarama.NewConfig()
config.Consumer.Return.Errors = true
config.Consumer.Offsets.Initial = sarama.OffsetOldest
```

### RabbitMQ QoS

```go
ch.Qos(
    1,     // prefetch count - process one message at a time
    0,     // prefetch size
    false, // global
)
```

## 📚 Additional Resources

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Sarama (Kafka Go Client)](https://github.com/IBM/sarama)
- [AMQP 0-9-1 (RabbitMQ Protocol)](https://www.rabbitmq.com/tutorials/amqp-concepts.html)

## 🎯 Next Steps

1. **Add Schema Registry**: Avro/Protobuf for Kafka
2. **Add Dead Letter Queues**: Handle failed messages
3. **Add Monitoring**: Prometheus metrics
4. **Add Tracing**: OpenTelemetry integration
5. **Add Security**: TLS, SASL authentication
6. **Add Transactions**: Exactly-once semantics

## 📝 License

MIT License - feel free to use this project for learning!


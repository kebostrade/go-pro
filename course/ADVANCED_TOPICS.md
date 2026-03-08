# Advanced Topics Learning Path

Specialized Go development tracks for mastering production-grade systems.

## Prerequisites

**Required**: Complete Lessons 1-10 (Foundations + Intermediate)  
**Recommended**: Complete Lessons 11-20 (Advanced + Expert)

## Track Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    ADVANCED TOPICS TRACK                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  TIER 1: Core Development (Prerequisite for all)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │ AT-01: REST │  │ AT-02: CLI  │  │ AT-03: Test │             │
│  │    APIs     │──│    Apps     │──│   & Debug   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
│  TIER 2: Framework & Architecture                               │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │ AT-04: Gin  │  │ AT-05: Micro│                              │
│  │  Framework  │──│  services   │                              │
│  └─────────────┘  └─────────────┘                              │
│                                                                 │
│  TIER 3: Communication Patterns                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │ AT-06: WS   │  │ AT-07: gRPC │  │ AT-08: NATS │             │
│  │  Real-time  │──│ Distributed │──│  Event-Drv  │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
│  TIER 4: Cloud & Infrastructure                                 │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │ AT-09: K8s  │  │ AT-10: Lambda│                              │
│  │ Cloud-Native│──│  Serverless │                              │
│  └─────────────┘  └─────────────┘                              │
│                                                                 │
│  TIER 5: Specialized Domains (Choose 1+)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │ AT-11: ML   │  │ AT-12: Block│  │ AT-13: IoT  │             │
│  │  Gorgonia   │  │  Ethereum   │  │    MQTT     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
│  TIER 6: API Design                                             │
│  ┌─────────────┐                                                │
│  │ AT-14: GraphQL│                                              │
│  │   gqlgen    │                                                │
│  └─────────────┘                                                │
│                                                                 │
│  CROSS-CUTTING: AT-00 System Design (with any tier)             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Detailed Learning Paths

### Path A: Backend API Developer
```
Lessons 1-10 → AT-01 REST APIs → AT-04 Gin → AT-14 GraphQL
             → AT-03 Testing → AT-05 Microservices → AT-07 gRPC
```
**Duration**: 6-8 weeks  
**Outcome**: Full-stack API development skills

### Path B: DevOps/Cloud Engineer
```
Lessons 1-10 → AT-02 CLI Apps → AT-05 Microservices
             → AT-09 Kubernetes → AT-10 Serverless → AT-08 NATS
```
**Duration**: 5-7 weeks  
**Outcome**: Cloud-native deployment expertise

### Path C: Distributed Systems Engineer
```
Lessons 1-10 → AT-00 System Design → AT-05 Microservices
             → AT-07 gRPC → AT-08 NATS → AT-09 Kubernetes
```
**Duration**: 6-8 weeks  
**Outcome**: Large-scale systems design

### Path D: Real-time/IoT Developer
```
Lessons 1-10 → AT-06 WebSockets → AT-08 NATS → AT-13 IoT
```
**Duration**: 4-5 weeks  
**Outcome**: Real-time systems expertise

### Path E: Blockchain/Web3 Developer
```
Lessons 1-10 → AT-01 REST APIs → AT-12 Blockchain → AT-14 GraphQL
```
**Duration**: 5-6 weeks  
**Outcome**: Web3 application development

---

## Topic Details

### TIER 1: Core Development

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-01 | RESTful APIs | 8-10h | Intermediate | L1-10 |
| AT-02 | CLI Applications | 6-8h | Intermediate | L1-10 |
| AT-03 | Testing & Debugging | 8-10h | Intermediate | L1-10 |

### TIER 2: Framework & Architecture

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-04 | Gin Web Framework | 8-10h | Intermediate | AT-01 |
| AT-05 | Microservices + Docker | 12-15h | Advanced | AT-01, AT-04 |

### TIER 3: Communication Patterns

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-06 | WebSockets Real-time | 8-10h | Advanced | AT-01 |
| AT-07 | gRPC Distributed | 10-12h | Advanced | AT-05 |
| AT-08 | NATS Event-Driven | 8-10h | Advanced | AT-05 |

### TIER 4: Cloud & Infrastructure

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-09 | Kubernetes Cloud-Native | 12-15h | Expert | AT-05 |
| AT-10 | AWS Lambda Serverless | 8-10h | Advanced | AT-01 |

### TIER 5: Specialized Domains

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-11 | ML with Gorgonia | 10-12h | Expert | L1-15 |
| AT-12 | Blockchain Ethereum | 12-15h | Expert | AT-01 |
| AT-13 | IoT with MQTT | 8-10h | Advanced | AT-06 |

### TIER 6: API Design

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-14 | GraphQL with gqlgen | 10-12h | Advanced | AT-01, AT-04 |

### CROSS-CUTTING

| ID | Topic | Duration | Difficulty | Prerequisites |
|----|-------|----------|------------|---------------|
| AT-00 | System Design | 10-12h | Expert | L1-15 |

---

## Recommended Schedule

### Full Track (All 15 Topics)
- **Weeks 1-2**: TIER 1 (AT-01, AT-02, AT-03)
- **Weeks 3-4**: TIER 2 (AT-04, AT-05)
- **Weeks 5-6**: TIER 3 (AT-06, AT-07, AT-08)
- **Weeks 7-8**: TIER 4 (AT-09, AT-10)
- **Weeks 9-10**: TIER 5 (Choose 2-3)
- **Week 11**: TIER 6 (AT-14)
- **Week 12**: AT-00 System Design + Capstone

### Quick Start (Essential 5)
1. AT-01: RESTful APIs
2. AT-03: Testing & Debugging
3. AT-04: Gin Framework
4. AT-05: Microservices
5. AT-00: System Design

---

## File Structure

```
course/
├── advanced-topics/
│   ├── AT-00-system-design/
│   │   ├── README.md
│   │   ├── examples/
│   │   ├── exercises/
│   │   └── projects/
│   ├── AT-01-rest-apis/
│   ├── AT-02-cli-apps/
│   ├── AT-03-testing-debugging/
│   ├── AT-04-gin-framework/
│   ├── AT-05-microservices-docker/
│   ├── AT-06-websockets/
│   ├── AT-07-grpc-distributed/
│   ├── AT-08-nats-event-driven/
│   ├── AT-09-kubernetes-cloud/
│   ├── AT-10-serverless-lambda/
│   ├── AT-11-ml-gorgonia/
│   ├── AT-12-blockchain-ethereum/
│   ├── AT-13-iot-mqtt/
│   └── AT-14-graphql-gqlgen/
└── projects/
    └── advanced/
        ├── rest-api-project/
        ├── microservices-project/
        ├── realtime-chat-project/
        ├── distributed-system-project/
        └── serverless-project/
```

---

## Assessment

Each topic includes:
- **Theory Quiz**: 10 questions
- **Coding Exercises**: 4-6 problems
- **Project**: 1 practical application
- **Code Review**: Peer review checklist

### Certification Tracks

1. **API Developer**: AT-01, AT-04, AT-14
2. **Backend Engineer**: AT-01, AT-03, AT-05, AT-07
3. **Cloud Engineer**: AT-05, AT-09, AT-10
4. **Full-Stack Go**: All Tier 1-4 + AT-14

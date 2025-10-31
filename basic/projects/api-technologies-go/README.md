# Tutorial 21: RESTful APIs, gRPC & GraphQL with Go

A comprehensive guide to building modern APIs with Go, covering three major API paradigms: REST, gRPC, and GraphQL.

## 📚 Overview

This tutorial teaches you how to build production-ready APIs using three different technologies:

1. **RESTful APIs** - HTTP-based APIs with JSON (using net/http, Chi, and Gin)
2. **gRPC** - High-performance RPC framework with Protocol Buffers
3. **GraphQL** - Flexible query language for APIs

## 🎯 Learning Objectives

By completing this tutorial, you will learn:

- Build RESTful APIs with different frameworks (net/http, Chi, Gin)
- Implement middleware for logging, CORS, authentication
- Create gRPC services with Protocol Buffers
- Implement streaming RPC (server, client, bidirectional)
- Build GraphQL servers with schema and resolvers
- Compare and choose the right API technology
- Build an API Gateway combining all three technologies

## 📋 Prerequisites

- Go 1.22 or higher
- Basic understanding of HTTP and APIs
- Familiarity with JSON
- Docker (optional, for infrastructure)

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd basic/projects/api-technologies-go
make deps
```

### 2. Run Examples

```bash
# REST API examples
make run-rest-basic    # Basic REST with net/http
make run-rest-chi      # REST with Chi router
make run-rest-gin      # REST with Gin framework

# gRPC examples
make proto-gen         # Generate protobuf code (first time only)
make run-grpc-unary    # Unary RPC
make run-grpc-streaming # Streaming RPC

# GraphQL example
make run-graphql       # GraphQL server

# Combined example
make run-gateway       # API Gateway (all three)
```

## 📁 Project Structure

```
api-technologies-go/
├── rest/
│   ├── basic/          # Basic REST API with net/http
│   ├── chi/            # REST API with Chi router
│   └── gin/            # REST API with Gin framework
├── grpc/
│   ├── proto/          # Protocol Buffer definitions
│   ├── unary/          # Unary RPC examples
│   └── streaming/      # Streaming RPC examples
├── graphql/
│   ├── schema/         # GraphQL schema definitions
│   └── gqlgen/         # GraphQL server
├── combined/
│   └── gateway/        # API Gateway (REST + gRPC + GraphQL)
├── docker-compose.yml  # Infrastructure services
├── Makefile           # Build automation
└── README.md          # This file
```

## 🔵 RESTful APIs

### Basic REST API (net/http)

Location: `rest/basic/main.go`

Features:
- Pure Go standard library
- Custom middleware (logging, CORS, JSON)
- CRUD operations
- Error handling
- In-memory storage

**Endpoints:**
```
GET    /users      - List all users
POST   /users      - Create new user
GET    /users/{id} - Get user by ID
PUT    /users/{id} - Update user
DELETE /users/{id} - Delete user
```

**Example Request:**
```bash
# List users
curl http://localhost:8080/users

# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com"}'
```

### Chi Router REST API

Location: `rest/chi/main.go`

Features:
- Chi router with middleware
- Request context
- Input validation with go-playground/validator
- CORS support
- Structured routing

**Endpoints:**
```
GET    /health              - Health check
GET    /api/v1/users        - List users
POST   /api/v1/users        - Create user
GET    /api/v1/users/{id}   - Get user
PUT    /api/v1/users/{id}   - Update user
DELETE /api/v1/users/{id}   - Delete user
```

### Gin Framework REST API

Location: `rest/gin/main.go`

Features:
- Gin web framework
- Built-in validation
- Query parameter filtering
- JSON binding
- Route grouping

**Endpoints:**
```
GET    /health                      - Health check
POST   /api/v1/login                - Login
GET    /api/v1/users                - List users (with filters)
POST   /api/v1/users                - Create user
GET    /api/v1/users/:id            - Get user
PUT    /api/v1/users/:id            - Update user
DELETE /api/v1/users/:id            - Delete user
PATCH  /api/v1/users/:id/activate   - Activate user
PATCH  /api/v1/users/:id/deactivate - Deactivate user
GET    /api/v1/admin/stats          - Get statistics
```

**Example with Filters:**
```bash
# Filter by role and active status
curl "http://localhost:8080/api/v1/users?role=admin&active=true"
```

## 🟣 gRPC

### Protocol Buffers

Location: `grpc/proto/user.proto`

Define your service and messages:
```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc GetUser(GetUserRequest) returns (User);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc StreamUsers(StreamUsersRequest) returns (stream User);
  rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

### Generate Code

```bash
make proto-gen
```

This generates:
- `user.pb.go` - Message types
- `user_grpc.pb.go` - Service definitions

### Unary RPC

Location: `grpc/unary/main.go`

Features:
- Simple request-response
- Connection pooling
- Error handling with status codes
- Context with timeout

**Example:**
```go
client := pb.NewUserServiceClient(conn)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := client.CreateUser(ctx, &pb.CreateUserRequest{
    Username: "john",
    Email:    "john@example.com",
    Role:     "user",
})
```

### Streaming RPC

Location: `grpc/streaming/main.go`

Features:
- Server streaming (one request, multiple responses)
- Client streaming (multiple requests, one response)
- Bidirectional streaming (both directions)

**Server Streaming Example:**
```go
stream, err := client.StreamUsers(ctx, &pb.StreamUsersRequest{})
for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Process user
}
```

**Client Streaming Example:**
```go
stream, err := client.CreateUsers(ctx)
for _, req := range requests {
    stream.Send(req)
}
resp, err := stream.CloseAndRecv()
```

**Bidirectional Streaming Example:**
```go
stream, err := client.Chat(ctx)

// Send messages
go func() {
    for _, msg := range messages {
        stream.Send(msg)
    }
    stream.CloseSend()
}()

// Receive messages
for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Process message
}
```

## 🟢 GraphQL

### Schema Definition

Location: `graphql/schema/schema.graphql`

```graphql
type User {
  id: ID!
  username: String!
  email: String!
  role: Role!
  posts: [Post!]!
}

type Query {
  user(id: ID!): User
  users(role: Role, active: Boolean): [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
}
```

### GraphQL Server

Location: `graphql/gqlgen/main.go`

Features:
- Schema-first development
- Type-safe resolvers
- Built-in playground
- Query validation

**Example Queries:**
```graphql
# Get all users
query GetUsers {
  users {
    id
    username
    email
    role
  }
}

# Get user with posts
query GetUser {
  user(id: "1") {
    id
    username
    posts {
      id
      title
    }
  }
}

# Create user
mutation CreateUser {
  createUser(input: {
    username: "john"
    email: "john@example.com"
    role: USER
  }) {
    id
    username
  }
}
```

## 🌐 API Gateway

Location: `combined/gateway/main.go`

A unified gateway that combines REST, gRPC, and GraphQL:

**Features:**
- Single entry point for all APIs
- Interactive dashboard
- API comparison guide
- Unified data store

**Access:**
```
http://localhost:8082          - Interactive dashboard
http://localhost:8082/api/rest/users      - REST endpoint
http://localhost:8082/api/graphql         - GraphQL endpoint
http://localhost:8082/api/grpc/users      - gRPC proxy
http://localhost:8082/api/comparison      - API comparison
```

## 📊 API Comparison

| Feature | REST | gRPC | GraphQL |
|---------|------|------|---------|
| **Protocol** | HTTP/1.1 | HTTP/2 | HTTP/1.1 |
| **Data Format** | JSON | Protocol Buffers | JSON |
| **Schema** | Optional (OpenAPI) | Required (.proto) | Required (.graphql) |
| **Streaming** | Limited (SSE) | Built-in | Subscriptions |
| **Performance** | Good | Excellent | Good |
| **Browser Support** | Excellent | Limited | Excellent |
| **Learning Curve** | Easy | Moderate | Moderate |
| **Caching** | Easy | Complex | Complex |

### When to Use Each

**REST:**
- Public APIs
- Simple CRUD operations
- Mobile apps
- Wide client compatibility needed

**gRPC:**
- Microservices communication
- Real-time streaming
- Internal APIs
- High-performance requirements

**GraphQL:**
- Complex data requirements
- Mobile apps (reduce bandwidth)
- Aggregating multiple data sources
- Rapid frontend development

## 🧪 Testing

```bash
# Run all tests
make test

# Test specific example
cd rest/basic && go test -v
cd grpc/unary && go test -v
cd graphql/gqlgen && go test -v
```

## 🔧 Development

### Hot Reload

Use Air for hot reload during development:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd rest/gin && air
```

### Linting

```bash
make lint
```

## 📚 Additional Resources

### REST
- [Go net/http documentation](https://pkg.go.dev/net/http)
- [Chi router](https://github.com/go-chi/chi)
- [Gin framework](https://github.com/gin-gonic/gin)

### gRPC
- [gRPC Go documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [gRPC best practices](https://grpc.io/docs/guides/performance/)

### GraphQL
- [gqlgen documentation](https://gqlgen.com/)
- [GraphQL specification](https://spec.graphql.org/)
- [GraphQL best practices](https://graphql.org/learn/best-practices/)

## 🎓 Learning Outcomes

After completing this tutorial, you will be able to:

✅ Build RESTful APIs with multiple frameworks  
✅ Implement middleware for cross-cutting concerns  
✅ Create gRPC services with Protocol Buffers  
✅ Implement all types of RPC streaming  
✅ Build GraphQL servers with schema and resolvers  
✅ Choose the right API technology for your use case  
✅ Build API gateways combining multiple technologies  
✅ Implement proper error handling and validation  
✅ Apply best practices for each API type  

## 🤝 Contributing

Found an issue or want to improve the tutorial? Contributions are welcome!

## 📝 License

This tutorial is part of the Go Pro Learning Platform.


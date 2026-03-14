# gRPC Distributed System Examples

Comprehensive examples of building distributed systems with gRPC in Go, covering all major patterns and best practices.

## Overview

This directory contains practical examples of:

- **Unary RPC**: Simple request-response pattern
- **Server Streaming**: Server sends multiple responses to one request
- **Client Streaming**: Client sends multiple requests to one response
- **Bidirectional Streaming**: Both sides send messages asynchronously
- **Interceptors**: Middleware for authentication, logging, etc.
- **Deadlines & Timeouts**: Managing request lifecycle
- **Error Handling**: Proper gRPC error status codes
- **TLS/SSL**: Secure communication
- **Reflection**: Runtime service inspection
- **Load Balancing**: Client-side load balancing

## Prerequisites

Install Protocol Buffer compiler:

```bash
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# Or download from https://github.com/protocolbuffers/protobuf/releases
```

Install Go plugins:

```bash
make deps
```

## Project Structure

```
08-grpc-distributed/
├── proto/                   # Protocol buffer definitions
│   ├── user_service.proto   # User service definition
│   ├── *.pb.go              # Generated Go code
│   └── Makefile             # Proto compilation
├── server/                  # gRPC server implementations
│   ├── main.go              # Server entry point
│   ├── server.go            # Server setup and configuration
│   ├── interceptors/        # Middleware
│   └── tls/                 # TLS configuration
├── client/                  # gRPC client implementations
│   ├── main.go              # Client entry point
│   ├── client.go            # Client setup and configuration
│   └── load_balancer.go     # Load balancing examples
├── examples/                # Individual pattern examples
│   ├── 01-unary-rpc.go
│   ├── 02-server-stream.go
│   ├── 03-client-stream.go
│   ├── 04-bidirectional-stream.go
│   └── 05-interceptors.go
└── README.md                # This file
```

## Quick Start

### 1. Generate Go Code from Proto Files

```bash
make proto
# or
make gen
```

This compiles `.proto` files and generates Go code in the `proto/` directory.

### 2. Run the Server

```bash
# Terminal 1
make server
# or
cd server && go run main.go
```

### 3. Run the Client

```bash
# Terminal 2
make client
# or
cd client && go run main.go
```

### 4. Run Individual Examples

```bash
# Unary RPC
cd examples && go run 01-unary-rpc.go

# Server streaming
cd examples && go run 02-server-stream.go

# Client streaming
cd examples && go run 03-client-stream.go

# Bidirectional streaming
cd examples && go run 04-bidirectional-stream.go
```

## Protocol Buffer Definitions

### user_service.proto

```protobuf
syntax = "proto3";

package user;

option go_package = "./proto";

// User service with all RPC patterns
service UserService {
  // Unary RPC
  rpc GetUser(GetUserRequest) returns (GetUserResponse);

  // Server streaming
  rpc ListUsers(ListUsersRequest) returns (stream User);

  // Client streaming
  rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);

  // Bidirectional streaming
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  int32 limit = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  int32 age = 3;
}

message CreateUsersResponse {
  repeated User users = 1;
  int32 count = 2;
}

message ChatMessage {
  string user_id = 1;
  string message = 2;
  int64 timestamp = 3;
}
```

## gRPC Patterns

### 1. Unary RPC

Simple request-response pattern. Client sends one request, server sends one response.

**Example:** `examples/01-unary-rpc.go`

```go
// Client
req := &user.GetUserRequest{Id: "123"}
resp, err := client.GetUser(ctx, req)

// Server
func (s *server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    // Business logic
    return &user.GetUserResponse{User: foundUser}, nil
}
```

### 2. Server Streaming

Server sends multiple responses to one request.

**Example:** `examples/02-server-stream.go`

```go
// Client
stream, err := client.ListUsers(ctx, &user.ListUsersRequest{Limit: 10})
for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Process user
}

// Server
func (s *server) ListUsers(req *user.ListUsersRequest, stream user.UserService_ListUsersServer) error {
    for _, user := range users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. Client Streaming

Client sends multiple requests, server sends one response.

**Example:** `examples/03-client-stream.go`

```go
// Client
stream, err := client.CreateUsers(ctx)
for _, user := range users {
    stream.Send(&user.CreateUserRequest{Name: user.Name})
}
resp, err := stream.CloseAndRecv()

// Server
func (s *server) CreateUsers(stream user.UserService_CreateUsersServer) error {
    var createdUsers []*user.User
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&user.CreateUsersResponse{
                Users: createdUsers,
                Count: int32(len(createdUsers)),
            })
        }
        // Process request
    }
}
```

### 4. Bidirectional Streaming

Both client and server send messages asynchronously.

**Example:** `examples/04-bidirectional-stream.go`

```go
// Client
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

// Server
func (s *server) Chat(stream user.UserService_ChatServer) error {
    // Start goroutine to receive messages
    go func() {
        for {
            req, err := stream.Recv()
            if err == io.EOF {
                return
            }
            // Process and send response
            stream.Send(response)
        }
    }()
    // Keep stream open
    <-ctx.Done()
    return nil
}
```

## Interceptors (Middleware)

### Server Interceptors

**Example:** `server/interceptors/`

```go
// Logging interceptor
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("%s took %v", info.FullMethod, time.Since(start))
    return resp, err
}

// Authentication interceptor
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    token := metadata.ValueFromIncomingContext(ctx, "authorization")
    if len(token) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }
    // Validate token
    return handler(ctx, req)
}
```

### Client Interceptors

```go
// Retry interceptor
func RetryInterceptor(maxRetries int) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        var err error
        for i := 0; i < maxRetries; i++ {
            err = invoker(ctx, method, req, reply, cc, opts...)
            if err == nil {
                return nil
            }
            time.Sleep(time.Second * time.Duration(i+1))
        }
        return err
    }
}
```

## Deadlines and Timeouts

```go
// Set deadline
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

resp, err := client.GetUser(ctx, req)
if err == context.DeadlineExceeded {
    log.Println("Request timed out")
}

// Server can check deadline
if dl, ok := ctx.Deadline(); ok {
    log.Printf("Deadline: %v", dl)
}
```

## Error Handling

Use proper gRPC status codes:

```go
import "google.golang.org/grpc/status"

// Server
if user == nil {
    return nil, status.Error(codes.NotFound, "user not found")
}

if req.Email == "" {
    return nil, status.Error(codes.InvalidArgument, "email is required")
}

// Client
resp, err := client.GetUser(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            // Handle not found
        case codes.InvalidArgument:
            // Handle invalid argument
        }
    }
}
```

## TLS/SSL Configuration

### Generate Certificates

```bash
make certs
```

### Server with TLS

```go
// server/main.go
creds, err := credentials.LoadTLSCredentials("certs/server.crt", "certs/server.key")
if err != nil {
    log.Fatal(err)
}

server := grpc.NewServer(
    grpc.Creds(creds),
)
```

### Client with TLS

```go
// client/main.go
creds := credentials.NewTLS(&tls.Config{
    InsecureSkipVerify: true, // Only for development!
})

conn, err := grpc.Dial("localhost:50051",
    grpc.WithTransportCredentials(creds),
)
```

## Reflection for Debugging

Enable reflection on server:

```go
import "google.golang.org/grpc/reflection"

// In server setup
reflection.Register(server)
```

Use grpcurl to inspect services:

```bash
# List all services
make reflect
# or
grpcurl -plaintext localhost:50051 list

# Describe a service
make describe
# or
grpcurl -plaintext localhost:50051 describe UserService

# Call a method
grpcurl -plaintext -d '{"id": "123"}' localhost:50051 user.UserService/GetUser
```

## Load Balancing

### Client-Side Load Balancing

```go
// Round-robin load balancing
conn, err := grpc.Dial(
    "localhost:50051,localhost:50052,localhost:50053",
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

## Make Commands

```bash
make help           # Show all available commands
make gen            # Generate Go code from .proto files
make server         # Run gRPC server
make client         # Run gRPC client
make examples       # Run all examples
make test           # Run tests
make lint           # Run linter
make fmt            # Format code
make clean          # Clean generated files
make deps           # Install dependencies
make certs          # Generate TLS certificates
make grpcurl        # Install grpcurl
make reflect        # List all services
make describe       # Describe UserService
```

## Testing

Run all tests:

```bash
make test
```

Run specific test:

```bash
cd server && go test -v -run TestGetUser
```

## Best Practices

### 1. Always Use Context

```go
// Good
func (s *server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Process request
    }
}
```

### 2. Handle Stream Errors

```go
for {
    user, err := stream.Recv()
    if err != nil {
        if err == io.EOF {
            break // Normal close
        }
        return err // Actual error
    }
    // Process user
}
```

### 3. Set Appropriate Deadlines

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 4. Use Metadata for Authentication

```go
// Client
md := metadata.Pairs("authorization", "Bearer "+token)
ctx := metadata.NewOutgoingContext(context.Background(), md)

// Server
token := metadata.ValueFromIncomingContext(ctx, "authorization")
```

### 5. Graceful Shutdown

```go
// Server
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
go func() {
    <-c
    server.GracefulStop()
}()
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 50051
lsof -i :50051

# Kill process
kill -9 <PID>
```

### Connection Refused

Ensure server is running:

```bash
make server
```

### Proto Files Not Generated

```bash
# Ensure protoc is installed
protoc --version

# Ensure plugins are installed
make deps

# Regenerate
make clean && make gen
```

## Additional Resources

- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [gRPC Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)

## License

MIT

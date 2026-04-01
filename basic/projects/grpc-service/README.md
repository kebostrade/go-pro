# gRPC Service Template

A gRPC distributed systems template demonstrating Protocol Buffers, all RPC streaming patterns, and proper error handling.

## gRPC Streaming Patterns

| RPC Type | Pattern | Use Case |
|----------|---------|----------|
| **Unary** | Request → Response | Simple queries |
| **Server Streaming** | Request → stream Responses | Large data sets |
| **Client Streaming** | stream Requests → Response | Uploads, batch operations |
| **Bidirectional** | stream ↔ stream | Real-time chat |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        gRPC Client                          │
│              (cmd/client/main.go examples)                   │
└──────────────────────────┬──────────────────────────────────┘
                           │ :50051
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     gRPC Server                              │
│                  (cmd/server/main.go)                       │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              UserService Implementation                 │  │
│  │   - GetUser (unary)                                   │  │
│  │   - ListUsers (server streaming)                      │  │
│  │   - CreateUsers (client streaming)                    │  │
│  │   - Chat (bidirectional streaming)                    │  │
│  └────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Generate protobuf code
make generate

# Start server
make server

# In another terminal, run client examples
make client
```

## Examples

### 1. Unary RPC (GetUser)
```go
resp, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: "1"})
```

### 2. Server Streaming (ListUsers)
```go
stream, err := client.ListUsers(ctx, &userpb.ListUsersRequest{Limit: 5})
for {
    user, err := stream.Recv()
    if err == io.EOF { break }
    // Process user
}
```

### 3. Client Streaming (CreateUsers)
```go
stream, err := client.CreateUsers(ctx)
stream.Send(&userpb.CreateUserRequest{Name: "Alice", ...})
stream.Send(&userpb.CreateUserRequest{Name: "Bob", ...})
resp, err := stream.CloseAndRecv()
```

### 4. Bidirectional Streaming (Chat)
```go
stream, err := client.Chat(ctx)
go func() {
    for {
        msg, _ := stream.Recv()
        // Handle incoming messages
    }
}()
// Send messages
stream.Send(&userpb.ChatMessage{Message: "Hello!"})
```

## Debugging with grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50051 list

# Describe UserService
grpcurl -plaintext localhost:50051 describe UserService

# Call GetUser
grpcurl -plaintext localhost:50051 UserService/GetUser -d '{"id": "1"}'

# Call ListUsers (server streaming)
grpcurl -plaintext localhost:50051 UserService/ListUsers -d '{"limit": 5}'
```

## Protocol Buffer Definition

```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc ListUsers(ListUsersRequest) returns (stream User);
  rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

## Environment Variables

- `GRPC_PORT` - Server port (default: 50051)

## Project Structure

```
grpc-service/
├── proto/
│   ├── user.proto              # Protocol buffer definitions
│   ├── user.pb.go              # Generated protobuf code
│   └── user_grpc.pb.go         # Generated gRPC code
├── cmd/
│   ├── server/main.go          # gRPC server
│   └── client/main.go          # Client examples
├── internal/
│   └── service/
│       └── service.go          # UserService implementation
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

## Docker

```bash
# Build
docker build -t grpc-service .

# Run
docker run -p 50051:50051 grpc-service
```

## License

MIT

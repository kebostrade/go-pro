# gRPC Distributed System Examples - Quick Start Guide

## Setup (One-Time)

### 1. Install Dependencies

```bash
# Install Protocol Buffer compiler
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. Generate Code from Proto Files

```bash
cd /home/dima/Desktop/FUN/go-pro/advanced-topics/08-grpc-distributed
make proto
```

Or use the setup script:

```bash
./setup.sh
```

### 3. Initialize Go Modules

```bash
go mod tidy
cd server && go mod tidy && cd ..
cd client && go mod tidy && cd ..
cd examples && go mod tidy && cd ..
```

## Run Examples

### Terminal 1: Start Server

```bash
cd server
go run main.go
```

**Output:**
```
🚀 Starting gRPC Server on port :50051
⚠️  Running without TLS (use 'make certs' to enable)
✓ Server registered and reflection enabled
📡 Listening on :50051
```

### Terminal 2: Run Client

```bash
cd client
go run main.go
```

**Output:**
```
🔌 Connecting to gRPC server at localhost:50051
✓ Connected to server

=== 📞 Unary RPC Example ===
✓ Received user: Alice Johnson

=== 📡 Server Streaming Example ===
✓ Received 3 users

=== 📤 Client Streaming Example ===
✓ Created 3 users

=== 💬 Bidirectional Streaming Example ===
✓ Received 5 messages

✓ All examples completed successfully!
```

## Run Individual Examples

### 1. Unary RPC

```bash
cd examples
go run 01-unary-rpc.go
```

**Demonstrates:**
- Simple request-response pattern
- Basic gRPC call
- Error handling

### 2. Server Streaming

```bash
cd examples
go run 02-server-stream.go
```

**Demonstrates:**
- Server sends multiple responses
- Client receives stream
- Stream termination (EOF)

### 3. Client Streaming

```bash
cd examples
go run 03-client-stream.go
```

**Demonstrates:**
- Client sends multiple requests
- Stream.Send() for each message
- Stream.CloseAndRecv() to get response

### 4. Bidirectional Streaming

```bash
cd examples
go run 04-bidirectional-stream.go
```

**Demonstrates:**
- Both sides send messages simultaneously
- Goroutines for concurrent send/receive
- Full-duplex communication

### 5. Interceptors

```bash
cd examples
go run 05-interceptors.go
```

**Demonstrates:**
- Logging interceptor
- Retry interceptor
- Authentication interceptor
- Chained interceptors

### 6. Deadlines

```bash
cd examples
go run 06-deadlines.go
```

**Demonstrates:**
- Context timeouts
- Deadline handling
- Request cancellation

### 7. Error Handling

```bash
cd examples
go run 07-error-handling.go
```

**Demonstrates:**
- gRPC status codes
- Error type checking
- Code-specific handling

### 8. Load Balancing

```bash
# First, start multiple server instances
# Terminal 1
cd server && go run main.go

# Terminal 2 (edit server/main.go to change port to 50052)
cd server && PORT=50052 go run main.go

# Terminal 3 (edit server/main.go to change port to 50053)
cd server && PORT=50053 go run main.go

# Terminal 4: Run load balancing example
cd examples
go run 08-load-balancing.go
```

**Demonstrates:**
- Manual round-robin
- gRPC automatic load balancing
- Connection pooling

### 9. TLS/SSL

```bash
# Generate certificates
make certs

# Server will automatically use TLS if certificates exist
cd server && go run main.go

# Run TLS example
cd examples
go run 09-tls.go
```

**Demonstrates:**
- Secure connections
- Certificate loading
- TLS configuration

### 10. Metadata

```bash
cd examples
go run 10-metadata.go
```

**Demonstrates:**
- Authentication tokens
- Tracing information
- Custom headers
- Multi-value metadata

## Debugging with grpcurl

### Install grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### List All Services

```bash
grpcurl -plaintext localhost:50051 list
```

**Output:**
```
grpc.reflection.v1.ServerReflection
user.UserService
```

### Describe Service

```bash
grpcurl -plaintext localhost:50051 describe user.UserService
```

### List Methods

```bash
grpcurl -plaintext localhost:50051 list user.UserService
```

**Output:**
```
user.UserService.GetUser
user.UserService.ListUsers
user.UserService.CreateUsers
user.UserService.Chat
```

### Call Method

```bash
grpcurl -plaintext \
  -d '{"id": "1"}' \
  localhost:50051 \
  user.UserService/GetUser
```

**Output:**
```json
{
  "user": {
    "id": "1",
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "age": 28,
    "tags": ["premium", "active"]
  }
}
```

## Common Issues

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
cd server && go run main.go
```

### Proto Files Not Generated

```bash
cd proto
make clean
make generate
```

### Import Errors

```bash
# Run from appropriate directory
cd server && go mod tidy
cd client && go mod tidy
cd examples && go mod tidy
```

## Testing

Run all tests:
```bash
make test
```

Run specific tests:
```bash
cd server && go test -v ./...
cd client && go test -v ./...
```

## Clean Up

Remove generated files:
```bash
make clean
```

Remove certificates:
```bash
rm -rf certs/
```

## Next Steps

1. **Experiment with proto files**: Add new services or messages
2. **Implement your own interceptors**: Add logging, metrics, etc.
3. **Try different streaming patterns**: Mix and match patterns
4. **Add error handling**: Implement custom status codes
5. **Explore production features**: Load balancing, service discovery, etc.

## Resources

- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [gRPC Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)

## Support

For issues or questions:
1. Check the main README.md
2. Review example code in examples/
3. Use grpcurl to inspect services
4. Check server logs for errors

# Building Distributed Systems with Go and gRPC

Design and implement high-performance distributed systems using gRPC.

## Learning Objectives

- Define protobuf messages
- Implement gRPC services
- Handle streaming RPCs
- Implement load balancing
- Add TLS security
- Debug gRPC services

## Theory

### Proto Definition

```protobuf
syntax = "proto3";

package userservice;

option go_package = "github.com/myorg/userservice/pb";

service UserService {
    rpc GetUser(GetUserRequest) returns (User) {}
    rpc ListUsers(ListUsersRequest) returns (stream User) {}
    rpc CreateUser(stream CreateUserRequest) returns (CreateUserResponse) {}
    rpc Chat(stream ChatMessage) returns (stream ChatMessage) {}
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    int64 created_at = 4;
}

message GetUserRequest {
    string id = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message CreateUserResponse {
    repeated User users = 1;
}

message ChatMessage {
    string message = 1;
}
```

### Server Implementation

```go
type userServiceServer struct {
    pb.UnimplementedUserServiceServer
    repo UserRepository
}

func NewUserServiceServer(repo UserRepository) *userServiceServer {
    return &userServiceServer{repo: repo}
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.repo.FindByID(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
    }

    return &pb.User{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt.Unix(),
    }, nil
}

func (s *userServiceServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users, err := s.repo.List(stream.Context(), req.PageSize, req.PageToken)
    if err != nil {
        return status.Errorf(codes.Internal, "list users: %v", err)
    }

    for _, user := range users {
        if err := stream.Send(&pb.User{
            Id:    user.ID,
            Name:  user.Name,
            Email: user.Email,
        }); err != nil {
            return err
        }
    }

    return nil
}

func (s *userServiceServer) Chat(stream pb.UserService_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        response := &pb.ChatMessage{
            Message: "Echo: " + msg.Message,
        }

        if err := stream.Send(response); err != nil {
            return err
        }
    }
}
```

### Server Setup

```go
func runServer() error {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        return fmt.Errorf("listen: %w", err)
    }

    opts := []grpc.ServerOption{
        grpc.MaxConcurrentStreams(1000),
        grpc.KeepaliveParams(keepalive.ServerParameters{
            MaxConnectionIdle: 5 * time.Minute,
            Time:              10 * time.Second,
            Timeout:           1 * time.Second,
        }),
    }

    if tlsCert, tlsKey := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"); tlsCert != "" {
        creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
        if err != nil {
            return fmt.Errorf("tls: %w", err)
        }
        opts = append(opts, grpc.Creds(creds))
    }

    s := grpc.NewServer(opts...)
    pb.RegisterUserServiceServer(s, NewUserServiceServer(NewUserRepo()))

    reflection.Register(s)

    return s.Serve(lis)
}
```

### Client Implementation

```go
type UserClient struct {
    client pb.UserServiceClient
    conn   *grpc.ClientConn
}

func NewUserClient(addr string) (*UserClient, error) {
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithDefaultCallOptions(
            grpc.MaxCallRecvMsgSize(10 * 1024 * 1024),
        ),
    }

    conn, err := grpc.Dial(addr, opts...)
    if err != nil {
        return nil, fmt.Errorf("dial: %w", err)
    }

    return &UserClient{
        client: pb.NewUserServiceClient(conn),
        conn:   conn,
    }, nil
}

func (c *UserClient) GetUser(ctx context.Context, id string) (*User, error) {
    resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
    if err != nil {
        return nil, err
    }

    return &User{
        ID:    resp.Id,
        Name:  resp.Name,
        Email: resp.Email,
    }, nil
}

func (c *UserClient) Close() error {
    return c.conn.Close()
}
```

## Security Considerations

```go
func withServerTLS(certFile, keyFile string) grpc.ServerOption {
    creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
    if err != nil {
        panic(err)
    }
    return grpc.Creds(creds)
}

func withClientTLS(caFile string) grpc.DialOption {
    creds, err := credentials.NewClientTLSFromFile(caFile, "")
    if err != nil {
        panic(err)
    }
    return grpc.WithTransportCredentials(creds)
}

func authInterceptor(authToken string) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        tokens := md.Get("authorization")
        if len(tokens) == 0 || tokens[0] != "Bearer "+authToken {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }

        return handler(ctx, req)
    }
}
```

## Performance Tips

```go
func withConnectionPool(addr string, poolSize int) ([]*grpc.ClientConn, error) {
    conns := make([]*grpc.ClientConn, poolSize)
    for i := 0; i < poolSize; i++ {
        conn, err := grpc.Dial(addr,
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithKeepaliveParams(keepalive.ClientParameters{
                Time:                10 * time.Second,
                Timeout:             time.Second,
                PermitWithoutStream: true,
            }),
        )
        if err != nil {
            return nil, err
        }
        conns[i] = conn
    }
    return conns, nil
}
```

## Exercises

1. Implement a key-value store service
2. Add bidirectional streaming
3. Implement client load balancing
4. Add TLS encryption

## Validation

```bash
cd exercises
protoc --go_out=. --go-grpc_out=. proto/*.proto
go test -v ./...
grpcurl -plaintext localhost:50051 list
```

## Key Takeaways

- Use protobuf for efficient serialization
- Choose appropriate RPC type (unary vs stream)
- Always implement error handling with status codes
- Use TLS in production
- Enable reflection for debugging

## Next Steps

**[AT-08: NATS Event-Driven](../AT-08-nats-event-driven/README.md)**

---

gRPC: fast, typed, streaming. ⚡

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

// Interceptors Example - Demonstrating client-side interceptors for logging, retry, and auth
func main() {
	fmt.Println("🔧 Interceptors Example")
	fmt.Println("======================")

	// Create connection with interceptors
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			loggingInterceptor,
			retryInterceptor,
			authInterceptor,
		),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Call RPC with interceptors
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("\n📤 Calling GetUser with interceptors...")
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Printf("\n✓ Received user: %s (%s)\n", resp.User.Name, resp.User.Email)
	fmt.Println("\n✅ Interceptors example completed!")
}

// Logging interceptor logs all RPC calls
func loggingInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	fmt.Printf("🔍 [LOG] Calling %s\n", method)

	err := invoker(ctx, method, req, reply, cc, opts...)

	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("❌ [LOG] %s failed: %v (took %v)\n", method, err, elapsed)
	} else {
		fmt.Printf("✓ [LOG] %s succeeded (took %v)\n", method, elapsed)
	}

	return err
}

// Retry interceptor retries failed requests
func retryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			if i > 0 {
				fmt.Printf("🔄 [RETRY] Success on attempt %d/%d\n", i+1, maxRetries)
			}
			return nil
		}

		lastErr = err
		fmt.Printf("⚠️  [RETRY] Attempt %d/%d failed: %v\n", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}
	}

	return lastErr
}

// Auth interceptor adds authentication metadata
func authInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// In a real application, you would add actual auth tokens
	// For demo purposes, we're just logging
	fmt.Println("🔐 [AUTH] Adding authentication metadata")

	// You can add metadata like this:
	// md := metadata.Pairs("authorization", "Bearer "+token)
	// ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

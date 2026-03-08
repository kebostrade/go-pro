package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/dima/go-pro/advanced-topics/08-grpc-distributed/proto"
)

// TLS/SSL Example - Demonstrating secure communication
//
// Before running this example:
// 1. Generate TLS certificates: make certs
// 2. Start server with TLS enabled
func main() {
	fmt.Println("🔒 TLS/SSL Example")
	fmt.Println("==================")

	// Example 1: Secure connection with TLS
	fmt.Println("\n📝 Example 1: Secure TLS connection")
	connectWithTLS()

	// Example 2: Insecure connection (for development)
	fmt.Println("\n📝 Example 2: Insecure connection")
	connectInsecure()

	// Example 3: TLS with custom verification
	fmt.Println("\n📝 Example 3: TLS with skip verify (development only)")
	connectWithSkipVerify()

	fmt.Println("\n✅ TLS example completed!")
	fmt.Println("\nNote: For TLS examples to work, generate certificates first:")
	fmt.Println("  make certs")
}

func connectWithTLS() {
	// Load TLS certificates
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		fmt.Printf("  ⚠️  Failed to load TLS certificates: %v\n", err)
		fmt.Println("  Run 'make certs' to generate certificates")
		return
	}

	// Create connection with TLS
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		fmt.Printf("  ❌ Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("  ✓ Connected with TLS")

	// Make RPC call
	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ RPC failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user over secure connection: %s\n", resp.User.Name)
}

func connectInsecure() {
	fmt.Println("  ⚠️  Using insecure connection (development only)")

	// For development/testing, you can use insecure credentials
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
	)
	if err != nil {
		fmt.Printf("  ❌ Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("  ✓ Connected with TLS (verification skipped)")

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ RPC failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s\n", resp.User.Name)
}

func connectWithSkipVerify() {
	fmt.Println("  ⚠️  Skipping certificate verification (NEVER use in production)")

	// Create custom TLS config that skips verification
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		fmt.Printf("  ❌ Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("  ✓ Connected (TLS verification disabled)")

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		fmt.Printf("  ❌ RPC failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Received user: %s\n", resp.User.Name)
}

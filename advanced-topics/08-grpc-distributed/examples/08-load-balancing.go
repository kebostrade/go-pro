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

// Load Balancing Example - Demonstrating client-side load balancing
//
// Note: This is a simplified example. In production, you would:
// 1. Run multiple server instances on different ports
// 2. Use a service discovery mechanism (DNS, etcd, Consul)
// 3. Implement proper round-robin or other load balancing strategies
func main() {
	fmt.Println("⚖️  Load Balancing Example")
	fmt.Println("==========================")

	// Simulate multiple server endpoints
	// In production, these would be actual running servers
	servers := []string{
		"localhost:50051",
		"localhost:50052",
		"localhost:50053",
	}

	// Example 1: Manual load balancing (round-robin)
	fmt.Println("\n📝 Example 1: Manual round-robin load balancing")
	manualRoundRobin(servers)

	// Example 2: gRPC built-in load balancing
	fmt.Println("\n📝 Example 2: gRPC automatic load balancing")
	grpcLoadBalancing()

	fmt.Println("\n✅ Load balancing example completed!")
	fmt.Println("\nNote: For full functionality, run multiple server instances")
	fmt.Println("on ports 50051, 50052, and 50053")
}

func manualRoundRobin(servers []string) {
	fmt.Println("  Creating connections to multiple servers...")

	conns := make([]*grpc.ClientConn, 0, len(servers))
	clients := make([]pb.UserServiceClient, 0, len(servers))

	// Connect to all servers
	for _, addr := range servers {
		fmt.Printf("  📞 Connecting to %s...\n", addr)

		conn, err := grpc.Dial(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			fmt.Printf("    ⚠️  Failed to connect: %v\n", err)
			continue
		}

		conns = append(conns, conn)
		clients = append(clients, pb.NewUserServiceClient(conn))
		fmt.Printf("    ✓ Connected\n")
	}

	if len(clients) == 0 {
		fmt.Println("  ❌ No servers available")
		return
	}

	// Distribute requests round-robin
	ctx := context.Background()
	requests := 10

	fmt.Printf("\n  Sending %d requests in round-robin fashion...\n", requests)

	for i := 0; i < requests; i++ {
		clientIndex := i % len(clients)
		client := clients[clientIndex]

		req := &pb.GetUserRequest{Id: "1"}
		resp, err := client.GetUser(ctx, req)

		if err != nil {
			fmt.Printf("  ❌ Request %d failed on server %d: %v\n", i+1, clientIndex+1, err)
			continue
		}

		fmt.Printf("  ✓ Request %d served by server %d: %s\n",
			i+1, clientIndex+1, resp.User.Name)

		time.Sleep(100 * time.Millisecond)
	}

	// Close connections
	for _, conn := range conns {
		conn.Close()
	}
}

func grpcLoadBalancing() {
	fmt.Println("  Using gRPC built-in load balancing...")

	// Configure gRPC with round-robin policy
	// Note: This requires multiple addresses or DNS resolution
	conn, err := grpc.Dial(
		"localhost:50051", // Single address for demo
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("  ❌ Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Make a few requests
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		req := &pb.GetUserRequest{Id: "1"}
		resp, err := client.GetUser(ctx, req)

		if err != nil {
			fmt.Printf("  ❌ Request %d failed: %v\n", i+1, err)
			continue
		}

		fmt.Printf("  ✓ Request %d: %s\n", i+1, resp.User.Name)
		time.Sleep(100 * time.Millisecond)
	}
}

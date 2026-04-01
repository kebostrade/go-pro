// Package main provides the gRPC server entry point.
package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/DimaJoyti/go-pro/basic/projects/grpc-service/internal/service"
	userpb "github.com/DimaJoyti/go-pro/basic/projects/grpc-service/proto"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	// Create listener.
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with options.
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(5 * time.Second),
	}
	grpcServer := grpc.NewServer(opts...)

	// Register user service.
	userServer := service.NewUserServiceServer()
	userpb.RegisterUserServiceServer(grpcServer, userServer)

	// Register reflection service for debugging with grpcurl.
	reflection.Register(grpcServer)

	// Register reflection service for debugging with grpcurl.
	reflection.Register(grpcServer)

	log.Printf("🚀 gRPC server starting on :%s", port)
	log.Printf("📋 Services: UserService")
	log.Printf("🔧 Debug: grpcurl -plaintext localhost:%s list", port)
	log.Printf("📖 Reflection: enabled")

	// Serve.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// GetPort returns the configured port.
func GetPort() string {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return "50051"
	}
	// Validate port.
	if p, err := strconv.Atoi(port); err != nil || p < 1 || p > 65535 {
		return "50051"
	}
	return port
}

// Ensure interface compliance.
var _ userpb.UserServiceServer = (*service.UserServiceServer)(nil)

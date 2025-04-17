package grpc

import (
	"fmt"
	"github/kijunpos/config"
	pbUser "github/kijunpos/gen/proto/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartGRPCServer starts the gRPC server
func StartGRPCServer(cfg *config.Config, userHandler UserHandler) {
	address := fmt.Sprintf(":%d", cfg.App.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pbUser.RegisterUserServiceServer(grpcServer, userHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

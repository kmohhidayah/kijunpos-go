package grpc

import (
	"github/kijunpos/app"
	"github/kijunpos/app/handler/user"
	pb "github/kijunpos/gen/proto/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(setup *app.SetupData) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Init user grpc
	pb.RegisterUserServiceServer(s, user.New(setup.ConfigData))

	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

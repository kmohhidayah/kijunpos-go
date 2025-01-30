package grpc

import (
	"github/kijunpos/app/handler"
	pb "github/kijunpos/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPosServiceServer(s, handler.NewPosHandler())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

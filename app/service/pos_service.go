package service

import (
	"context"
	pb "github/kijunpos/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosService struct {
}

func NewPosService() *PosService {
	return &PosService{}
}

func (s *PosService) GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (*pb.GetCustomerDetailResponse, error) {
	// Implementasi logika bisnis Anda di sini
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID cannot be empty")
	}
	now := time.Now()
	// Contoh response
	return &pb.GetCustomerDetailResponse{
		Id:          req.Id,
		Name:        req.Name,
		PhoneNumber: "123456789",          // Implement your logic here
		CreatedAt:   timestamppb.New(now), // Use time.Now().Proto() for actual timestamp
	}, nil
}

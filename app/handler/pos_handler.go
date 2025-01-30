package handler

import (
	"context"
	"github/kijunpos/app/service"
	"github/kijunpos/config"
	pb "github/kijunpos/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosHandler struct {
	pb.UnimplementedPosServiceServer
}

func NewPosHandler(cfg *config.Config, svc *service.PosService) *PosHandler {
	return &PosHandler{}
}

func (h *PosHandler) GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (*pb.GetCustomerDetailResponse, error) {

	// Business Logic
	id := req.Id
	return &pb.GetCustomerDetailResponse{
		Id:          id,
		Name:        "John Doe",
		PhoneNumber: "081143211234",
		CreatedAt:   timestamppb.New(time.Now()),
	}, nil

}

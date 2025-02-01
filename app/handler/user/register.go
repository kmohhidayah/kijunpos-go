package user

import (
	"context"
	pb "github/kijunpos/gen/proto/user"
)

func (h *grpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.GeneralResponse, error) {
	return h.UserService.Register(ctx)
}

package user

import (
	"context"
	"github/kijunpos/app/lib/apm"
	pb "github/kijunpos/gen/proto/user"
)

func (h *grpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.GeneralResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "handler.user.register")
	defer span.End()
	return h.UserService.Register(ctx)
}

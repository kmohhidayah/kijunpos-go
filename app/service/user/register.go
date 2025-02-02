package user

import (
	"context"
	"github/kijunpos/app/lib/apm"
	pbUser "github/kijunpos/gen/proto/user"
)

func (s *Service) Register(ctx context.Context) (*pbUser.GeneralResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "service.user.Register")
	defer span.End()

	return &pbUser.GeneralResponse{
		Message: "success di panggil dari layer service",
		Success: true,
	}, nil
}

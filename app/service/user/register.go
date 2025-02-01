package user

import (
	"context"
	pbUser "github/kijunpos/gen/proto/user"
)

func (s Service) Register(ctx context.Context) (*pbUser.GeneralResponse, error) {
	return &pbUser.GeneralResponse{
		Message: "success di panggil dari layer service",
		Success: true,
	}, nil
}

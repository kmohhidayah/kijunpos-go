package user

import (
	"github/kijunpos/app/service/user"
	"github/kijunpos/config"
	pbUser "github/kijunpos/gen/proto/user"
)

type grpcHandler struct {
	pbUser.UnimplementedUserServiceServer
	UserService user.Service
}

func New(cfg *config.Config) *grpcHandler {
	return &grpcHandler{
		UserService: user.New(cfg),
	}
}

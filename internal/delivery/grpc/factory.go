package grpc

import (
	pbUser "github/kijunpos/gen/proto/user"
	userHandler "github/kijunpos/internal/delivery/grpc/user"
	"github/kijunpos/internal/domain"
)

// UserHandler interface for gRPC user handler
type UserHandler interface {
	pbUser.UserServiceServer
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase domain.UserUseCase) UserHandler {
	return userHandler.NewHandler(userUseCase)
}

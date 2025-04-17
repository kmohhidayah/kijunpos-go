package user

import (
	"github/kijunpos/internal/domain"
	pbUser "github/kijunpos/gen/proto/user"
)

// Handler handles gRPC requests for user service
type Handler struct {
	pbUser.UnimplementedUserServiceServer
	userUseCase domain.UserUseCase
}

// NewHandler creates a new user handler
func NewHandler(userUseCase domain.UserUseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

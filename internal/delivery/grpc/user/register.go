package user

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	pbUser "github/kijunpos/gen/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register handles user registration
func (h *Handler) Register(ctx context.Context, req *pbUser.RegisterRequest) (*pbUser.GeneralResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.Register")
	defer span.End()

	// Validate request
	if req.Name == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "name, password, and email are required")
	}

	// Create user domain object
	user := &domain.User{
		UserName:     req.Name,
		PasswordHash: req.Password, // Will be hashed in the use case
		Email:        req.Email,
	}

	// Call use case
	_, err := h.userUseCase.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbUser.GeneralResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

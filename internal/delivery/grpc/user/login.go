package user

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	pbUser "github/kijunpos/gen/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login handles user authentication
func (h *Handler) Login(ctx context.Context, req *pbUser.LoginRequest) (*pbUser.LoginResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.Login")
	defer span.End()

	// Determine auth type based on request
	var authType domain.AuthType
	switch req.AuthType {
	case "email":
		authType = domain.AuthTypeEmail
	case "whatsapp":
		authType = domain.AuthTypeWhatsApp
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid auth type")
	}

	// Validate input
	if req.Identifier == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}
	if req.Credential == "" {
		return nil, status.Error(codes.InvalidArgument, "credential is required")
	}

	// Call use case
	user, err := h.userUseCase.Login(ctx, authType, req.Identifier, req.Credential)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Create response
	return &pbUser.LoginResponse{
		Success: true,
		Message: "Login successful",
		User: &pbUser.UserData{
			Id:          user.ID.String(),
			Username:    user.UserName,
			Email:       user.Email,
			PhoneNumber: user.WhatsAppNumber,
		},
	}, nil
}

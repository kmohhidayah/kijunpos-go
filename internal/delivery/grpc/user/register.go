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

	// Determine registration type based on request
	var authType domain.AuthType
	params := make(map[string]string)

	if req.PhoneNumber != "" {
		// Register with WhatsApp (using phone number as WhatsApp number)
		authType = domain.AuthTypeWhatsApp
		params["whatsapp_number"] = req.PhoneNumber
	} else if req.Email != "" && req.Password != "" {
		// Register with Email
		authType = domain.AuthTypeEmail
		params["email"] = req.Email
		params["password"] = req.Password
	} else {
		return nil, status.Error(codes.InvalidArgument, "either phone number or email and password are required")
	}

	// Validate username
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Call use case
	_, err := h.userUseCase.Register(ctx, authType, req.Name, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbUser.GeneralResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

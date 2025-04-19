package user

import (
	"context"
	pbUser "github/kijunpos/gen/proto/user"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"github/kijunpos/internal/pkg/errors"
)

// Login handles user authentication
func (h *Handler) Login(ctx context.Context, req *pbUser.LoginRequest) (*pbUser.GeneralResponse, error) {
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
		return &pbUser.GeneralResponse{
			Success: false,
			Message: errors.NewBadRequestError("invalid auth type", nil).Error(),
		}, nil
	}

	// Validate input
	if req.Identifier == "" {
		return &pbUser.GeneralResponse{
			Success: false,
			Message: errors.NewValidationError("identifier is required", nil).Error(),
		}, nil
	}
	if req.Credential == "" {
		return &pbUser.GeneralResponse{
			Success: false,
			Message: errors.NewValidationError("credential is required", nil).Error(),
		}, nil
	}

	// Call use case
	_, err := h.userUseCase.Login(ctx, authType, req.Identifier, req.Credential)
	if err != nil {
		// Handle the error using the custom error package
		errMsg := errors.HandleResponseError(ctx, span, err)
		return &pbUser.GeneralResponse{
			Success: false,
			Message: errMsg,
		}, nil
	}

	// Create response
	return &pbUser.GeneralResponse{
		Success: true,
		Message: "Login successful",
	}, nil
}

package user

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"github/kijunpos/internal/pkg/errors"
	pbUser "github/kijunpos/gen/proto/user"
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
		return errors.NewErrorResponse("Either phone number or email and password are required"), nil
	}

	// Validate username
	if req.Name == "" {
		return errors.NewErrorResponse("Name is required"), nil
	}

	// Call use case
	_, err := h.userUseCase.Register(ctx, authType, req.Name, params)
	if err != nil {
		return errors.MapErrorToResponse(ctx, span, err)
	}

	return errors.NewSuccessResponse("User registered successfully"), nil
}

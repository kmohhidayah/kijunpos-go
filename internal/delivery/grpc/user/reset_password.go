package user

import (
	"context"
	pbUser "github/kijunpos/gen/proto/user"
	"github/kijunpos/internal/pkg/apm"
	"github/kijunpos/internal/pkg/errors"
)

// ResetPassword handles password reset requests
func (h *Handler) ResetPassword(ctx context.Context, req *pbUser.ResetPasswordRequest) (*pbUser.ResetPasswordResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.ResetPassword")
	defer span.End()

	// Validate input
	if req.Email == "" {
		return &pbUser.ResetPasswordResponse{
			Success: false,
			Message: errors.NewValidationError("email is required", nil).Error(),
		}, nil
	}

	// Call use case
	verificationCode, err := h.userUseCase.ResetPassword(ctx, req.Email)
	if err != nil {
		errMsg := errors.HandleResponseError(ctx, span, err)
		return &pbUser.ResetPasswordResponse{
			Success: false,
			Message: errMsg,
		}, nil
	}

	// In a real application, the verification code would be sent to the user via email
	// and not returned in the response
	return &pbUser.ResetPasswordResponse{
		Success:          true,
		Message:          "Verification code sent to your email",
		VerificationCode: verificationCode, // Only for testing purposes
	}, nil
}

// VerifyPasswordReset handles password reset verification
func (h *Handler) VerifyPasswordReset(ctx context.Context, req *pbUser.VerifyPasswordResetRequest) (*pbUser.GeneralResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.VerifyPasswordReset")
	defer span.End()

	// Validate input
	if req.Email == "" {
		return errors.NewErrorResponse("email is required"), nil
	}
	if req.VerificationCode == "" {
		return errors.NewErrorResponse("verification code is required"), nil
	}
	if req.NewPassword == "" {
		return errors.NewErrorResponse("new password is required"), nil
	}
	if req.NewPassword != req.ConfirmationPassword {
		return errors.NewErrorResponse("passwords do not match"), nil
	}

	// Call use case
	err := h.userUseCase.VerifyPasswordReset(ctx, req.Email, req.VerificationCode, req.NewPassword)
	if err != nil {
		return errors.MapErrorToResponse(ctx, span, err)
	}

	return errors.NewSuccessResponse("Password has been reset successfully"), nil
}

package user

import (
	"context"
	"github/kijunpos/internal/pkg/apm"
	pbUser "github/kijunpos/gen/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResetPIN handles PIN reset requests
func (h *Handler) ResetPIN(ctx context.Context, req *pbUser.ResetPINRequest) (*pbUser.ResetPINResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.ResetPIN")
	defer span.End()

	// Validate input
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}

	// Call use case
	verificationCode, err := h.userUseCase.ResetPIN(ctx, req.PhoneNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// In a real application, the verification code would be sent to the user via WhatsApp
	// and not returned in the response
	return &pbUser.ResetPINResponse{
		Success:          true,
		Message:          "Verification code sent to your WhatsApp number",
		VerificationCode: verificationCode, // Only for testing purposes
	}, nil
}

// VerifyPINReset handles PIN reset verification
func (h *Handler) VerifyPINReset(ctx context.Context, req *pbUser.VerifyPINResetRequest) (*pbUser.GeneralResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.user.VerifyPINReset")
	defer span.End()

	// Validate input
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}
	if req.VerificationCode == "" {
		return nil, status.Error(codes.InvalidArgument, "verification code is required")
	}
	if req.NewPin == "" {
		return nil, status.Error(codes.InvalidArgument, "new PIN is required")
	}

	// Call use case
	err := h.userUseCase.VerifyPINReset(ctx, req.PhoneNumber, req.VerificationCode, req.NewPin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbUser.GeneralResponse{
		Success: true,
		Message: "PIN has been reset successfully",
	}, nil
}

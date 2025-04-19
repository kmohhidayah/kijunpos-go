package errors

import (
	"context"
	pbUser "github/kijunpos/gen/proto/user"
	"go.opentelemetry.io/otel/trace"
)

// NewSuccessResponse creates a new GeneralResponse with success status
func NewSuccessResponse(message string) *pbUser.GeneralResponse {
	return &pbUser.GeneralResponse{
		Success: true,
		Message: message,
	}
}

// NewErrorResponse creates a new GeneralResponse with error status
func NewErrorResponse(message string) *pbUser.GeneralResponse {
	return &pbUser.GeneralResponse{
		Success: false,
		Message: message,
	}
}

// HandleError creates an error response using the HandleResponseError function
func HandleError(ctx context.Context, span trace.Span, err error) *pbUser.GeneralResponse {
	message := HandleResponseError(ctx, span, err)
	return NewErrorResponse(message)
}

// MapErrorToResponse maps an error to a GeneralResponse
// This is a convenience function that can be used directly in handlers
func MapErrorToResponse(ctx context.Context, span trace.Span, err error) (*pbUser.GeneralResponse, error) {
	return HandleError(ctx, span, err), nil
}

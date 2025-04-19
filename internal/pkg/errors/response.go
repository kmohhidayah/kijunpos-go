package errors

import (
	"context"
	"errors"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// ErrorCategory defines the category of an error
type ErrorCategory int

const (
	// CategoryValidation represents validation errors
	CategoryValidation ErrorCategory = iota
	// CategoryNotFound represents not found errors
	CategoryNotFound
	// CategoryDuplicate represents duplicate/already exists errors
	CategoryDuplicate
	// CategoryUnauthorized represents unauthorized access errors
	CategoryUnauthorized
	// CategoryForbidden represents forbidden access errors
	CategoryForbidden
	// CategoryBadRequest represents bad request errors
	CategoryBadRequest
	// CategoryInternal represents internal server errors
	CategoryInternal
)

// AppError represents an application error with category
type AppError struct {
	Category ErrorCategory
	Message  string
	Err      error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) error {
	return &AppError{
		Category: CategoryValidation,
		Message:  message,
		Err:      err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) error {
	return &AppError{
		Category: CategoryNotFound,
		Message:  message,
		Err:      err,
	}
}

// NewDuplicateError creates a new duplicate error
func NewDuplicateError(message string, err error) error {
	return &AppError{
		Category: CategoryDuplicate,
		Message:  message,
		Err:      err,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, err error) error {
	return &AppError{
		Category: CategoryUnauthorized,
		Message:  message,
		Err:      err,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, err error) error {
	return &AppError{
		Category: CategoryForbidden,
		Message:  message,
		Err:      err,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, err error) error {
	return &AppError{
		Category: CategoryBadRequest,
		Message:  message,
		Err:      err,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, err error) error {
	return &AppError{
		Category: CategoryInternal,
		Message:  message,
		Err:      err,
	}
}

// HandleResponseError maps error messages to appropriate response messages
// and logs internal errors to the provided span
func HandleResponseError(ctx context.Context, span trace.Span, err error) string {
	if err == nil {
		return ""
	}

	// Try to unwrap as AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		// For internal errors, don't expose details
		if appErr.Category == CategoryInternal {
			if span != nil {
				span.RecordError(err)
			}
			return "Operation failed, please try again later"
		}
		// For other categories, return the message
		return appErr.Error()
	}

	// Legacy error handling for errors not wrapped in AppError
	errMsg := err.Error()
	
	// Try to categorize based on error message
	switch {
	case strings.Contains(errMsg, "already exists"):
		return errMsg
	case strings.Contains(errMsg, "not found"):
		return errMsg
	case strings.Contains(errMsg, "is required") || 
	     strings.Contains(errMsg, "invalid") || 
	     strings.Contains(errMsg, "must be"):
		return errMsg
	case strings.Contains(errMsg, "unauthorized") || 
	     strings.Contains(errMsg, "unauthenticated"):
		return errMsg
	case strings.Contains(errMsg, "forbidden") || 
	     strings.Contains(errMsg, "permission denied"):
		return errMsg
	default:
		// Log the internal error but don't expose details to client
		if span != nil {
			span.RecordError(err)
		}
		return "Operation failed, please try again later"
	}
}

package domain

import "context"

// EmailService represents the email service contract
type EmailService interface {
	// SendVerificationCode sends a verification code to the specified email address
	SendVerificationCode(ctx context.Context, email, code string) error
}

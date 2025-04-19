package domain

import (
	"context"
	"time"
)

// VerificationRepository represents the verification repository contract
type VerificationRepository interface {
	// StoreVerificationCode stores a verification code for the given email with an expiration time
	StoreVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error
	
	// GetVerificationCode retrieves a verification code for the given email
	GetVerificationCode(ctx context.Context, email string) (string, error)
	
	// DeleteVerificationCode deletes a verification code for the given email
	DeleteVerificationCode(ctx context.Context, email string) error
}

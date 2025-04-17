package user

import (
	"context"
	"errors"
	"fmt"
	"github/kijunpos/internal/pkg/apm"
	"math/rand"
	"time"
)

// ResetPIN initiates a PIN reset process for a user with the given WhatsApp number
// Returns a verification code that would be sent to the user via WhatsApp
func (uc *userUseCase) ResetPIN(ctx context.Context, whatsAppNumber string) (string, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.ResetPIN")
	defer span.End()

	// Validate input
	if whatsAppNumber == "" {
		return "", errors.New("whatsapp number is required")
	}

	// Get user by WhatsApp number
	user, err := uc.userRepo.GetByWhatsAppNumber(ctx, whatsAppNumber)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return "", errors.New("user account is not active")
	}

	// Generate verification code
	verificationCode := generateVerificationCode()

	// In a real application, you would:
	// 1. Store the verification code in a temporary storage with an expiration time
	// 2. Send the verification code to the user via WhatsApp
	// 3. Return success to the caller

	// For simplicity, we're just returning the verification code
	// In a real application, you would not return this to the caller
	return verificationCode, nil
}

// VerifyPINReset verifies the reset code and sets a new PIN for the user
func (uc *userUseCase) VerifyPINReset(ctx context.Context, whatsAppNumber, verificationCode, newPIN string) error {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.VerifyPINReset")
	defer span.End()

	// Validate input
	if whatsAppNumber == "" {
		return errors.New("whatsapp number is required")
	}
	if verificationCode == "" {
		return errors.New("verification code is required")
	}
	if newPIN == "" {
		return errors.New("new PIN is required")
	}

	// Get user by WhatsApp number
	user, err := uc.userRepo.GetByWhatsAppNumber(ctx, whatsAppNumber)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// In a real application, you would:
	// 1. Verify the verification code against the stored one
	// 2. Check if the verification code has expired
	// 3. If valid, update the user's PIN

	// For simplicity, we're just checking a hardcoded verification code
	// In a real application, you would verify against a stored code
	if verificationCode != "123456" { // This should be replaced with actual verification
		return errors.New("invalid verification code")
	}

	// Update the user's PIN
	user.PIN = newPIN
	now := time.Now()
	user.UpdatedAt.Time = now
	user.UpdatedAt.Valid = true
	user.PasswordChangedAt.Time = now
	user.PasswordChangedAt.Valid = true

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// generateVerificationCode generates a 6-digit verification code
func generateVerificationCode() string {
	// Initialize random number generator with current time
	rand.Seed(time.Now().UnixNano())
	
	// Generate a random 6-digit number
	code := rand.Intn(900000) + 100000
	
	return fmt.Sprintf("%d", code)
}

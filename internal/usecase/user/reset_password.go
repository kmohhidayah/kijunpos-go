package user

import (
	"context"
	"errors"
	"fmt"
	"github/kijunpos/internal/pkg/apm"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ResetPassword initiates a password reset process for a user with the given email
// Returns a verification code that would be sent to the user via email
func (uc *userUseCase) ResetPassword(ctx context.Context, email string) (string, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.ResetPassword")
	defer span.End()

	// Validate input
	if email == "" {
		return "", errors.New("email is required")
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
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

	// Store the verification code with 10 minutes expiration
	if err := uc.verificationRepo.StoreVerificationCode(ctx, email, verificationCode, 10*time.Minute); err != nil {
		return "", fmt.Errorf("failed to store verification code: %w", err)
	}

	// Send the verification code via email
	if err := uc.emailService.SendVerificationCode(ctx, email, verificationCode); err != nil {
		// If sending fails, delete the stored code to prevent inconsistency
		_ = uc.verificationRepo.DeleteVerificationCode(ctx, email)
		return "", fmt.Errorf("failed to send verification code: %w", err)
	}

	// In a production environment, we would not return the verification code
	// But for testing purposes, we'll return it
	return verificationCode, nil
}

// VerifyPasswordReset verifies the reset code and sets a new password for the user
func (uc *userUseCase) VerifyPasswordReset(ctx context.Context, email, verificationCode, newPassword string) error {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.VerifyPasswordReset")
	defer span.End()

	// Validate input
	if email == "" {
		return errors.New("email is required")
	}
	if verificationCode == "" {
		return errors.New("verification code is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify the verification code
	storedCode, err := uc.verificationRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	if storedCode != verificationCode {
		return errors.New("invalid verification code")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	user.PasswordHash = string(hashedPassword)
	now := time.Now()
	user.UpdatedAt.Time = now
	user.UpdatedAt.Valid = true
	user.PasswordChangedAt.Time = now
	user.PasswordChangedAt.Valid = true

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete the verification code after successful password reset
	if err := uc.verificationRepo.DeleteVerificationCode(ctx, email); err != nil {
		// Just log the error, don't fail the operation
		// In a real application, you would log this error
		fmt.Printf("Failed to delete verification code: %v\n", err)
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

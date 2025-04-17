package user

import (
	"context"
	"database/sql"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user based on auth type
func (uc *userUseCase) Login(ctx context.Context, authType domain.AuthType, identifier, credential string) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.Login")
	defer span.End()

	// Validate input
	if identifier == "" {
		return nil, errors.New("identifier is required")
	}
	if credential == "" {
		return nil, errors.New("credential is required")
	}

	var user *domain.User
	var err error

	switch authType {
	case domain.AuthTypeEmail:
		// Try to get user by username first
		user, err = uc.userRepo.GetByUsername(ctx, identifier)
		if err != nil {
			return nil, err
		}

		// If user not found by username, try email
		if user == nil {
			user, err = uc.userRepo.GetByEmail(ctx, identifier)
			if err != nil {
				return nil, err
			}
		}

		if user == nil {
			return nil, errors.New("invalid username/email or password")
		}

		// Check if user is active
		if !user.IsActive {
			return nil, errors.New("user account is not active")
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credential))
		if err != nil {
			// Increment failed login attempts
			user.FailedLoginAttempts++
			_ = uc.userRepo.Update(ctx, user) // Ignore error for simplicity
			return nil, errors.New("invalid username/email or password")
		}

	case domain.AuthTypeWhatsApp:
		// Get user by WhatsApp number
		user, err = uc.userRepo.GetByWhatsAppNumber(ctx, identifier)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("invalid WhatsApp number or PIN")
		}

		// Check if user is active
		if !user.IsActive {
			return nil, errors.New("user account is not active")
		}

		// Verify PIN
		if user.PIN != credential {
			// Increment failed login attempts
			user.FailedLoginAttempts++
			_ = uc.userRepo.Update(ctx, user) // Ignore error for simplicity
			return nil, errors.New("invalid WhatsApp number or PIN")
		}

	default:
		return nil, errors.New("invalid auth type")
	}

	// Update last login time and reset failed login attempts
	now := time.Now()
	user.LastLoginAt = sql.NullTime{Time: now, Valid: true}
	user.FailedLoginAttempts = 0
	user.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Clear sensitive data before returning
	user.PasswordHash = ""
	user.PIN = ""
	return user, nil
}

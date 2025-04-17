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

// Login authenticates a user
func (uc *userUseCase) Login(ctx context.Context, username, password string) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.Login")
	defer span.End()

	// Get user by username
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		_ = uc.userRepo.Update(ctx, user) // Ignore error for simplicity
		return nil, errors.New("invalid username or password")
	}

	// Update last login time and reset failed login attempts
	now := time.Now()
	user.LastLoginAt = sql.NullTime{Time: now, Valid: true}
	user.FailedLoginAttempts = 0
	user.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

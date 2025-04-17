package user

import (
	"context"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register registers a new user
func (uc *userUseCase) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.Register")
	defer span.End()

	// Check if user with the same email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if user with the same username already exists
	existingUser, err = uc.userRepo.GetByUsername(ctx, user.UserName)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Prepare user data
	user.ID = uuid.New()
	user.PasswordHash = string(hashedPassword)
	user.IsActive = true
	user.FailedLoginAttempts = 0
	user.CreatedAt = time.Now()

	// Create user in the database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

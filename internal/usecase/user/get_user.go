package user

import (
	"context"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"

	"github.com/google/uuid"
)

// GetUserByID retrieves a user by ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.GetUserByID")
	defer span.End()

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

package user

import (
	"context"
	"errors"
	"github/kijunpos/internal/pkg/apm"

	"github.com/google/uuid"
)

// DeleteUser deletes a user
func (uc *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.DeleteUser")
	defer span.End()

	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return uc.userRepo.Delete(ctx, id)
}

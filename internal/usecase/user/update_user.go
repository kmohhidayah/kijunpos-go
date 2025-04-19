package user

import (
	"context"
	"database/sql"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"time"
)

// UpdateUser updates a user
func (uc *userUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.UpdateUser")
	defer span.End()

	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Update timestamp
	now := time.Now()
	user.UpdatedAt = sql.NullTime{Time: now, Valid: true}

	return uc.userRepo.Update(ctx, user)
}

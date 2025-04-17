package user

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
)

// Update updates a user in the database
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.Update")
	defer span.End()

	query := `
		UPDATE users
		SET 
			username = $2,
			password_hash = $3,
			email = $4,
			is_active = $5,
			failed_login_attempts = $6,
			last_login_at = $7,
			password_changed_at = $8,
			updated_at = $9
		WHERE id = $1
	`

	_, err := r.dbConn.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.UserName,
		user.PasswordHash,
		user.Email,
		user.IsActive,
		user.FailedLoginAttempts,
		user.LastLoginAt,
		user.PasswordChangedAt,
		user.UpdatedAt,
	)

	return err
}

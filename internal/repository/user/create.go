package user

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
)

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.Create")
	defer span.End()

	query := `
		INSERT INTO users (
			id, username, password_hash, email, is_active, 
			failed_login_attempts, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
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
		user.CreatedAt,
	)

	return err
}

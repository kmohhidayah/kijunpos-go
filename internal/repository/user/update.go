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
			whatsapp_number = $5,
			pin = $6,
			is_active = $7,
			failed_login_attempts = $8,
			last_login_at = $9,
			password_changed_at = $10,
			updated_at = $11
		WHERE id = $1
	`

	_, err := r.dbConn.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.UserName,
		user.PasswordHash,
		user.Email,
		user.WhatsAppNumber,
		user.PIN,
		user.IsActive,
		user.FailedLoginAttempts,
		user.LastLoginAt,
		user.PasswordChangedAt,
		user.UpdatedAt,
	)

	return err
}

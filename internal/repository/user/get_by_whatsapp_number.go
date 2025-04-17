package user

import (
	"context"
	"database/sql"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
)

// GetByWhatsAppNumber retrieves a user by their WhatsApp number
func (r *userRepository) GetByWhatsAppNumber(ctx context.Context, whatsAppNumber string) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.GetByWhatsAppNumber")
	defer span.End()

	query := `
		SELECT id, username, password_hash, email, whatsapp_number, pin, is_active, 
		failed_login_attempts, created_at, last_login_at, password_changed_at, updated_at, deleted_at
		FROM users
		WHERE whatsapp_number = $1 AND deleted_at IS NULL
	`

	var user domain.User
	var lastLoginAt, passwordChangedAt, updatedAt, deletedAt sql.NullTime

	err := r.dbConn.DB.QueryRowContext(ctx, query, whatsAppNumber).Scan(
		&user.ID,
		&user.UserName,
		&user.PasswordHash,
		&user.Email,
		&user.WhatsAppNumber,
		&user.PIN,
		&user.IsActive,
		&user.FailedLoginAttempts,
		&user.CreatedAt,
		&lastLoginAt,
		&passwordChangedAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.LastLoginAt = lastLoginAt
	user.PasswordChangedAt = passwordChangedAt
	user.UpdatedAt = updatedAt
	user.DeletedAt = deletedAt

	return &user, nil
}

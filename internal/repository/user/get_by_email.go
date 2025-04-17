package user

import (
	"context"
	"database/sql"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
)

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.GetByEmail")
	defer span.End()

	query := `
		SELECT 
			id, username, password_hash, email, is_active, 
			failed_login_attempts, created_at, last_login_at, 
			password_changed_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user domain.User
	err := r.dbConn.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

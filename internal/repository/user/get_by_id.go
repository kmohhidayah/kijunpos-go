package user

import (
	"context"
	"database/sql"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"

	"github.com/google/uuid"
)

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.GetByID")
	defer span.End()

	query := `
		SELECT 
			id, username, password_hash, email, is_active, 
			failed_login_attempts, created_at, last_login_at, 
			password_changed_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user domain.User
	err := r.dbConn.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

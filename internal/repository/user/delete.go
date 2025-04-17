package user

import (
	"context"
	"github/kijunpos/internal/pkg/apm"

	"github.com/google/uuid"
)

// Delete soft deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.Delete")
	defer span.End()

	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := r.dbConn.DB.ExecContext(ctx, query, id)
	return err
}

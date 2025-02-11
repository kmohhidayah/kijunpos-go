package kijundb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID    `db:"id"`
	UserName            string       `db:"username"`
	PasswordHash        string       `db:"password_hash"`
	Email               string       `db:"email"`
	IsActive            bool         `db:"is_active"`
	FailedLoginAttempts int          `db:"failed_login_attempts"`
	CreatedAt           time.Time    `db:"created_at"`
	LastLoginAt         sql.NullTime `db:"last_login_at"`
	PasswordChangedAt   sql.NullTime `db:"password_changed_at"`
	UpdatedAt           sql.NullTime `db:"updated_at"`
	DeletedAt           sql.NullTime `db:"deleted_at"`
}

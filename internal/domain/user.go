package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity
type User struct {
	ID                  uuid.UUID    `db:"id"`
	UserName            string       `db:"username"`
	PasswordHash        string       `db:"password_hash"`
	Email               string       `db:"email"`
	WhatsAppNumber      string       `db:"whatsapp_number"`
	OTPPIN              string       `db:"otp_pin"`
	IsActive            bool         `db:"is_active"`
	FailedLoginAttempts int          `db:"failed_login_attempts"`
	CreatedAt           time.Time    `db:"created_at"`
	LastLoginAt         sql.NullTime `db:"last_login_at"`
	PasswordChangedAt   sql.NullTime `db:"password_changed_at"`
	UpdatedAt           sql.NullTime `db:"updated_at"`
	DeletedAt           sql.NullTime `db:"deleted_at"`
}

// AuthType defines the type of authenticatiuon (login or registration)
type AuthType string

const (
	// AuthTypeEmail represents authentication using email/username and password
	AuthTypeEmail AuthType = "email"
	// AuthTypeWhatsApp represents authentication using WhatsApp number and PIN
	AuthTypeWhatsApp AuthType = "whatsapp"
)

// UserRepository represents the user repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByWhatsAppNumber(ctx context.Context, whatsAppNumber string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserUseCase represents the user use case contract
type UserUseCase interface {
	Register(ctx context.Context, authType AuthType, username string, params map[string]string) (*User, error)
	Login(ctx context.Context, authType AuthType, identifier, credential string) (*User, error)
	ResetPassword(ctx context.Context, email string) (string, error)
	VerifyPasswordReset(ctx context.Context, email, verificationCode, newPassword string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

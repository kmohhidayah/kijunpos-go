package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity
type User struct {
	ID                  uuid.UUID    
	UserName            string       
	PasswordHash        string       
	Email               string       
	WhatsAppNumber      string       
	PIN                 string       
	IsActive            bool         
	FailedLoginAttempts int          
	CreatedAt           time.Time    
	LastLoginAt         sql.NullTime 
	PasswordChangedAt   sql.NullTime 
	UpdatedAt           sql.NullTime 
	DeletedAt           sql.NullTime 
}

// AuthType defines the type of authentication (login or registration)
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
	ResetPIN(ctx context.Context, whatsAppNumber string) (string, error)
	VerifyPINReset(ctx context.Context, whatsAppNumber, verificationCode, newPIN string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

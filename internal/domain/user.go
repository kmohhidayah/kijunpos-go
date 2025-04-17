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
	IsActive            bool         
	FailedLoginAttempts int          
	CreatedAt           time.Time    
	LastLoginAt         sql.NullTime 
	PasswordChangedAt   sql.NullTime 
	UpdatedAt           sql.NullTime 
	DeletedAt           sql.NullTime 
}

// UserRepository represents the user repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserUseCase represents the user use case contract
type UserUseCase interface {
	Register(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, username, password string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

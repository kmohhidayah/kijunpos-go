package user

import (
	"github/kijunpos/config/db"
	"github/kijunpos/internal/domain"
)

type userRepository struct {
	dbConn *db.Connection
}

// NewUserRepository creates a new user repository
func NewUserRepository(dbConn *db.Connection) domain.UserRepository {
	return &userRepository{
		dbConn: dbConn,
	}
}

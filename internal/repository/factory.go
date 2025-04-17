package repository

import (
	"github/kijunpos/config/db"
	"github/kijunpos/internal/domain"
	userRepo "github/kijunpos/internal/repository/user"
)

// NewUserRepository creates a new user repository
func NewUserRepository(dbConn *db.Connection) domain.UserRepository {
	return userRepo.NewUserRepository(dbConn)
}

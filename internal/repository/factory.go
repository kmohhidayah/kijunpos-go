package repository

import (
	"github/kijunpos/config/db"
	"github/kijunpos/internal/domain"
	userRepo "github/kijunpos/internal/repository/user"
	verificationRepo "github/kijunpos/internal/repository/verification"
)

// NewUserRepository creates a new user repository
func NewUserRepository(dbConn *db.Connection) domain.UserRepository {
	return userRepo.NewUserRepository(dbConn)
}

// NewVerificationRepository creates a new verification repository
func NewVerificationRepository() domain.VerificationRepository {
	return verificationRepo.NewVerificationRepository()
}

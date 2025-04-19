package user

import (
	"github/kijunpos/internal/domain"
)

type userUseCase struct {
	userRepo             domain.UserRepository
	verificationRepo     domain.VerificationRepository
	emailService         domain.EmailService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo domain.UserRepository,
	verificationRepo domain.VerificationRepository,
	emailService domain.EmailService,
) domain.UserUseCase {
	return &userUseCase{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		emailService:     emailService,
	}
}

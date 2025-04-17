package user

import (
	"context"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register registers a new user based on registration type
func (uc *userUseCase) Register(ctx context.Context, authType domain.AuthType, username string, params map[string]string) (*domain.User, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.user.Register")
	defer span.End()

	// Validate common input
	if username == "" {
		return nil, errors.New("username is required")
	}

	// Check if user with the same username already exists
	existingUser, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Create new user with common fields
	user := &domain.User{
		ID:                  uuid.New(),
		UserName:            username,
		IsActive:            true,
		FailedLoginAttempts: 0,
		CreatedAt:           time.Now(),
	}

	switch authType {
	case domain.AuthTypeWhatsApp:
		// Get WhatsApp number from params
		whatsAppNumber, ok := params["whatsapp_number"]
		if !ok || whatsAppNumber == "" {
			return nil, errors.New("whatsapp number is required")
		}

		// Check if user with the same WhatsApp number already exists
		existingUser, err = uc.userRepo.GetByWhatsAppNumber(ctx, whatsAppNumber)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("user with this WhatsApp number already exists")
		}

		// Generate a temporary PIN for the user
		// In a real application, this would be sent to the user via WhatsApp
		temporaryPIN := generateTemporaryPIN()

		// Set WhatsApp-specific fields
		user.WhatsAppNumber = whatsAppNumber
		user.PIN = temporaryPIN

	case domain.AuthTypeEmail:
		// Get email and password from params
		email, ok := params["email"]
		if !ok || email == "" {
			return nil, errors.New("email is required")
		}
		password, ok := params["password"]
		if !ok || password == "" {
			return nil, errors.New("password is required")
		}

		// Check if user with the same email already exists
		existingUser, err = uc.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("user with this email already exists")
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		// Set email-specific fields
		user.Email = email
		user.PasswordHash = string(hashedPassword)

	default:
		return nil, errors.New("invalid registration type")
	}

	// Create user in the database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Clear sensitive data before returning
	user.PasswordHash = ""
	user.PIN = ""
	return user, nil
}

// generateTemporaryPIN generates a 6-digit PIN for WhatsApp authentication
// In a real application, this would be more secure and sent to the user via WhatsApp
func generateTemporaryPIN() string {
	// For simplicity, we're using a fixed PIN here
	// In a real application, you would generate a random PIN
	return "123456"
}

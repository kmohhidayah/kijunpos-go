package verification

import (
	"context"
	"errors"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"sync"
	"time"
)

// Repository implements the domain.VerificationRepository interface
type Repository struct {
	codes      map[string]codeInfo
	mutex      sync.RWMutex
	cleanupJob *time.Ticker
}

type codeInfo struct {
	code      string
	expiresAt time.Time
}

// NewVerificationRepository creates a new verification repository
func NewVerificationRepository() domain.VerificationRepository {
	repo := &Repository{
		codes:      make(map[string]codeInfo),
		cleanupJob: time.NewTicker(5 * time.Minute),
	}

	// Start cleanup goroutine
	go func() {
		for range repo.cleanupJob.C {
			repo.cleanup()
		}
	}()

	return repo
}

// StoreVerificationCode stores a verification code for the given email with an expiration time
func (r *Repository) StoreVerificationCode(ctx context.Context, email, code string, expiration time.Duration) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.verification.StoreVerificationCode")
	defer span.End()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.codes[email] = codeInfo{
		code:      code,
		expiresAt: time.Now().Add(expiration),
	}

	return nil
}

// GetVerificationCode retrieves a verification code for the given email
func (r *Repository) GetVerificationCode(ctx context.Context, email string) (string, error) {
	ctx, span := apm.GetTracer().Start(ctx, "repository.verification.GetVerificationCode")
	defer span.End()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	info, exists := r.codes[email]
	if !exists {
		return "", errors.New("verification code not found")
	}

	if time.Now().After(info.expiresAt) {
		return "", errors.New("verification code has expired")
	}

	return info.code, nil
}

// DeleteVerificationCode deletes a verification code for the given email
func (r *Repository) DeleteVerificationCode(ctx context.Context, email string) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.verification.DeleteVerificationCode")
	defer span.End()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.codes, email)
	return nil
}

// cleanup removes expired codes
func (r *Repository) cleanup() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for email, info := range r.codes {
		if now.After(info.expiresAt) {
			delete(r.codes, email)
		}
	}
}

// Close stops the cleanup job
func (r *Repository) Close() {
	r.cleanupJob.Stop()
}

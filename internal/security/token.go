package security

import (
	"time"

	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/models"
)

const (
	TokenScopeAccess  = "access"
	TokenScopeRefresh = "refresh"
)

// Maker makes a new token
type Maker interface {

	// CreateToken creates a new token for a specific username and duration
	CreateToken(userID uuid.UUID, duration time.Duration, scope string) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)

	// GenerateAuthAccessToken generates a token at login
	GenerateAuthAccessToken(user *models.User, payload *AuthPayload, duration time.Duration) error
}

type AuthPayload struct {
	Data map[string]any
}

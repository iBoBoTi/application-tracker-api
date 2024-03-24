package security

import (
	"github.com/google/uuid"
	"time"
)

const (
	TokenScopeAccess  = "access"
	TokenScopeRefresh = "refresh"
)

// Maker makes a new token
type Maker interface {

	// CreateToken creates a new token for a specific username and duration
	CreateToken(userID uuid.UUID, duration time.Duration, version int64, scope string) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}

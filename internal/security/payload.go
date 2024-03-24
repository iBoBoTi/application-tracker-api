package security

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Different types of error that returned from the VerifyToken
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	IssuedAt  time.Time
	ExpiredAt time.Time
	Version   int64
	Scope     string
}

// NewPayload creates a new token with a specific username and duration
func NewPayload(userID uuid.UUID, duration time.Duration, scope string) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Scope:     scope,
	}

	return payload, nil
}

func (p *Payload) Valid() error {

	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

package security

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/o1egl/paseto"

	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/models"
)

// PasetoMaker is a Paseto token maker
type PasetoMaker struct {
	peseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: key must be %d characters", chacha20poly1305.KeySize)
	}

	maker := PasetoMaker{
		peseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return &maker, nil
}

func (m *PasetoMaker) GenerateAuthAccessToken(user *models.User, payload *AuthPayload, duration time.Duration) error {

	accessToken, accessTokenPayload, err := m.CreateToken(user.ID, duration, TokenScopeAccess)
	if err != nil {
		return appError.ErrInternalServer
	}
	payload.Data["accessToken"] = map[string]any{
		"token":     accessToken,
		"expiresAt": accessTokenPayload.ExpiredAt,
	}
	return nil
}

// CreateToken creates a new token for a specific username and duration
func (m *PasetoMaker) CreateToken(userID uuid.UUID, duration time.Duration, scope string) (string, *Payload, error) {

	payload, err := NewPayload(userID, duration, scope)
	if err != nil {
		return "", nil, err
	}
	tokenString, err := m.peseto.Encrypt(m.symmetricKey, payload, nil)

	return tokenString, payload, err
}

// VerifyToken checks if the token is valid or not
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}
	err := m.peseto.Decrypt(token, m.symmetricKey, payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

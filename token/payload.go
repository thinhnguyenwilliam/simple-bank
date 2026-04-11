// simple-bank/token/payload.go
package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDuration  = errors.New("invalid token duration")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSecretKey = errors.New("invalid secret key")
	ErrInvalidToken     = errors.New("invalid token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"` //token ID
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, role string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID:        uuid.New(),
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

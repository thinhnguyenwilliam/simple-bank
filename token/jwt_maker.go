// simple-bank/token/jwt_maker.go
package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf(
			"%w: must be at least %d characters",
			ErrInvalidSecretKey,
			minSecretKeySize,
		)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

func (m *JWTMaker) CreateToken(username, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, "test_role", duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         payload.ID.String(),
		"username":   payload.Username,
		"issued_at":  payload.IssuedAt.Unix(),
		"expired_at": payload.ExpiredAt.Unix(),
	})

	token, err := jwtToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", payload, err
	}

	return token, payload, nil
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// prevent algorithm confusion attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        id,
		Username:  claims["username"].(string),
		IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
		ExpiredAt: time.Unix(int64(claims["expired_at"].(float64)), 0),
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}

// simple-bank/token/jwt_maker_test.go
package token

import (
	"simple-bank/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	_, err := NewJWTMaker("short")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidSecretKey)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotNil(t, maker)
}

func TestJWTMakerCreateAndVerifyToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
}

func TestJWTMakerExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrExpiredToken)
}

func TestJWTMakerInvalidToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid.token.value")
	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTMakerWrongSecretKey(t *testing.T) {
	maker1, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	maker2, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker1.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token)
	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTMakerInvalidAlgNone(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	// Create token with alg = none
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id":         uuid.New().String(),
		"username":   "thinh",
		"issued_at":  time.Now().Unix(),
		"expired_at": time.Now().Add(time.Minute).Unix(),
	})

	// VERY IMPORTANT: allow unsigned token creation (test only)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(tokenString)

	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrInvalidToken)
}

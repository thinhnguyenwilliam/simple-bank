// simple-bank\token\paseto_maker_test.go

package token

import (
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMakerCreateAndVerifyToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := "thinh"
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
}

func TestPasetoMakerExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken("thinh", -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)

	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrExpiredToken)
}

func TestPasetoMakerWrongSecretKey(t *testing.T) {
	maker1, _ := NewPasetoMaker(util.RandomString(32))
	maker2, _ := NewPasetoMaker(util.RandomString(32))

	token, err := maker1.CreateToken("thinh", time.Minute)
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token)

	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrInvalidToken)
}

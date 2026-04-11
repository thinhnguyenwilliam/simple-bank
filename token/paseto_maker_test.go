// simple-bank\token\paseto_maker_test.go

package token

import (
	"testing"
	"time"

	"github.com/thinhcompany/simple-bank/util"

	"github.com/stretchr/testify/require"
)

func TestPasetoMakerCreateAndVerifyToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := "thinh"
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, "test_role", duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
}

func TestPasetoMakerExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken("thinh", "test_role", -time.Minute)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)

	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrExpiredToken)
}

func TestPasetoMakerWrongSecretKey(t *testing.T) {
	maker1, _ := NewPasetoMaker(util.RandomString(32))
	maker2, _ := NewPasetoMaker(util.RandomString(32))

	token, payload, err := maker1.CreateToken("thinh", "test_role", time.Minute)
	require.NoError(t, err)

	payload, err = maker2.VerifyToken(token)

	require.Nil(t, payload)
	require.ErrorIs(t, err, ErrInvalidToken)
}

// simple-bank\util\password_test.go

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	// hashed password must not equal plain password
	require.NotEqual(t, password, hashedPassword)
}

func TestCheckPassword(t *testing.T) {
	password := RandomString(10)
	wrongPassword := RandomString(10)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// correct password → no error
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	// wrong password → error
	err = CheckPassword(wrongPassword, hashedPassword)
	require.Error(t, err)
	require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}

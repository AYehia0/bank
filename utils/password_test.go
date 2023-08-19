package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	// check if the two passwords are equal
	password := RandomString(8)

	hashedPassword, err := GenerateHash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ComparePasswords(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := "WrongAssPassword"
	err = ComparePasswords(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// returns the hased password
func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("Failed to generate the hash for : %s", password)
	}

	return string(hashedPassword), nil
}

func ComparePasswords(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

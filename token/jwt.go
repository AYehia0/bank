// implements the token creator interface

package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const MIN_SECRET_LEN = 32

type JWTCreator struct {
	secretkey string
}

func NewJWTCreator(secret string) (TokenCreator, error) {
	if len(secret) < MIN_SECRET_LEN {
		return nil, fmt.Errorf("Invalid key length: must be at least %d", MIN_SECRET_LEN)
	}
	return &JWTCreator{
		secretkey: secret,
	}, nil
}

// create the token for the user with a specific time
func (c *JWTCreator) Create(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// make the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(c.secretkey))
}

func (c *JWTCreator) Verify(token string) (*Payload, error) {

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// we expect the token to be signed with HMAC algo only like jwt.SigningMethodHS256
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, TokenInvalidError
		}
		return []byte(c.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, TokenExpiredError) {
			return nil, TokenExpiredError
		}
		return nil, TokenInvalidError
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, TokenInvalidError
	}
	return payload, nil
}

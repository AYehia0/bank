package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	TokenExpiredError = errors.New("Token has been expired!")
	TokenInvalidError = errors.New("Token is invalid!")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Id:        tokenId,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// the Valid func is required to validate the payload (not verificaton) : required by jwt
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return TokenExpiredError
	}
	return nil
}

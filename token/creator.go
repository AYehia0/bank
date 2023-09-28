package token

import "time"

type TokenCreator interface {
	// create a token for a username/email with a duration time
	Create(username string, duration time.Duration) (string, *Payload, error)

	// verify the token
	Verify(token string) (*Payload, error)
}

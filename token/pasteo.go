package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var (
	InvalidKeySize = fmt.Errorf("Invalid key size: key must be %d", chacha20poly1305.KeySize)
)

type PasteoCreator struct {
	versionImp *paseto.V2

	// symmetric key
	secretKey []byte
}

func NewPasteoToken(secretKey string) (TokenCreator, error) {
	// the symmetric key used by pasteo is chachakey with a specific length, so the length of the key must statsfiy that length
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, InvalidKeySize
	}
	creator := &PasteoCreator{
		versionImp: paseto.NewV2(),
		secretKey:  []byte(secretKey),
	}

	return creator, nil
}

func (p *PasteoCreator) Create(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := p.versionImp.Encrypt(p.secretKey, payload, nil)
	return token, payload, err
}

func (p *PasteoCreator) Verify(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.versionImp.Decrypt(token, p.secretKey, payload, nil)

	if err != nil {
		return nil, TokenInvalidError
	}

	// validate the token : expired ?
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

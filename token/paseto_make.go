package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

var _ Maker = PasetoMake{}

// PasetoMake is a PASETO token maker
type PasetoMake struct {
	paseto *paseto.V2
	// symmetricKey 密钥
	symmetricKey []byte
}

func (p PasetoMake) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	accessToken, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	return accessToken, payload, err
}

func (p PasetoMake) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return PasetoMake{symmetricKey: []byte(symmetricKey), paseto: paseto.NewV2()}, nil
}

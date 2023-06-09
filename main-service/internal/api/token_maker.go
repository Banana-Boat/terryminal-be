package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type TokenPayload struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type TokenMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewTokenMaker(symmetricKey string) (*TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &TokenMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *TokenMaker) CreateToken(id int64, email string, duration time.Duration) (string, error) {
	payload := &TokenPayload{
		ID:        id,
		Email:     email,
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, err
}

func (maker *TokenMaker) VerifyToken(token string) (*TokenPayload, error) {
	payload := &TokenPayload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.New("token is invalid")
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, errors.New("token has expired")
	}

	return payload, nil
}

package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of the errors returned by VerfiyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	UserType  string    `json:"user_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expired_at"`
}

func NewPayload(tokenParams *TokenParams) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    tokenParams.UserID,
		Email:     tokenParams.Email,
		UserType:  tokenParams.UserType,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(tokenParams.Duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

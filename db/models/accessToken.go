package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type AccessToken interface{}

type accessToken struct {
	ID             string    `json:"id,omitempty"`
	createdAt      time.Time `json:"created_at"`
	expiresAt      time.Time `json:"expires_at"`
	refreshTokenID string    `json:"refresh_token_id,omitempty"`
}

func NewAccessToken(refreshTokenID string) (AccessToken, string) {
	ID := uuid.New().String()
	hash := sha256.Sum256([]byte(ID))

	return &accessToken{
		ID:             hex.EncodeToString(hash[:]),
		createdAt:      time.Now(),
		expiresAt:      time.Now().Add(time.Minute * 15),
		refreshTokenID: refreshTokenID,
	}, ID
}

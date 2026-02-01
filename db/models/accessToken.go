package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type AccessToken interface{}

// Access Token: token (hash), created_at, expires_at, fk => refresh_token
type accessToken struct {
	ID             string
	createdAt      time.Time
	expiresAt      time.Time
	refreshTokenID string
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

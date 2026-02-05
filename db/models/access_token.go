package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID             string    `json:"id,omitempty"`
	createdAt      time.Time `json:"created_at"`
	expiresAt      time.Time `json:"expires_at"`
	refreshTokenID string    `json:"refresh_token_id,omitempty"`
}

func NewAccessToken(userID uuid.UUID) (*AccessToken, string) {
	//ID := uuid.New().String()
	//hash := sha256.Sum256([]byte(ID))

	// Create Access token
	// Flow
	// User logs in
	// Create refreshToken + accessToken
	// Store refreshToken in DB - is it though?
	// Access token invalidates itself through expiration
	// Refresh token is created and send to the user..
	// Why is this stateful.. it has expiration jsut like the accessToken ?
	// Stateful so the server can revoke the refreshToken and invalidate it on logout ex.

	return &AccessToken{
		ID:             hex.EncodeToString(hash[:]),
		createdAt:      time.Now(),
		expiresAt:      time.Now().Add(time.Minute * 15),
		refreshTokenID: refreshTokenID,
	}, ID
}

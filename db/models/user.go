package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type User interface {
}

// User: id, name, email, role, refresh_token (hash), refresh_created_at
type user struct {
	id                    uuid.UUID `json:"id,omitempty"`
	name                  string    `json:"name,omitempty"`
	email                 string    `json:"email,omitempty"`
	role                  string    `json:"role,omitempty"` // TODO: Create enum
	refreshToken          string    `json:"refresh_token,omitempty"`
	refreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

func NewUser(name, email, role string) User {
	return &user{
		id:    uuid.New(),
		name:  name,
		email: email,
		role:  role,
	}
}

func (u *user) CreateRefreshToken() string {
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	u.refreshToken = hex.EncodeToString(hash[:])

	return token
}

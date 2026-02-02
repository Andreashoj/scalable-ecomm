package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type User interface {
	CreateRefreshToken() string
}

type user struct {
	id                    uuid.UUID `json:"id,omitempty"`
	name                  string    `json:"name,omitempty"`
	password              string    `json:"password,omitempty"`
	email                 string    `json:"email,omitempty"`
	role                  string    `json:"role,omitempty"` // TODO: Create enum
	refreshToken          string    `json:"refresh_token,omitempty"`
	refreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

func NewUser(name, email, password, role string) User {
	hashedPass := sha256.Sum256([]byte(password))

	user := &user{
		id:       uuid.New(),
		name:     name,
		email:    email,
		password: hex.EncodeToString(hashedPass[:]), // TODO: hash it
		role:     role,
	}

	user.CreateRefreshToken()
	return user
}

func (u *user) CreateRefreshToken() string {
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	u.refreshToken = hex.EncodeToString(hash[:])
	u.refreshTokenExpiresAt = time.Now().Add(time.Minute * 15)

	return token
}

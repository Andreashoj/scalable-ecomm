package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       uuid.UUID `json:"id,omitempty"`
	name     string    `json:"name,omitempty"`
	password string    `json:"password,omitempty"`
	email    string    `json:"email,omitempty"`
	role     Role      `json:"role,omitempty"` // TODO: Save enum
	//refreshToken          string    `json:"refresh_token,omitempty"` // move away from the user
	//refreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

type Role string

const (
	Customer Role = "customer"
)

func NewUser(name, email, password string) (*User, error) {
	hashedPass, err := createHashedPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed creating hashed pass: %s", err)
	}

	user := &User{
		id:       uuid.New(),
		name:     name,
		email:    email,
		password: hashedPass,
		role:     Customer,
	}

	return user, nil
}

func (u *User) CreateRefreshToken() string {
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	u.refreshToken = hex.EncodeToString(hash[:])
	u.refreshTokenExpiresAt = time.Now().Add(time.Minute * 15)

	return token
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

func createHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generating hashed password: %s", err)
	}

	return string(hashedPass), nil
}

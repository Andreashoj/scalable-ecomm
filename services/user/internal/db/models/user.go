package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      Role      `json:"role,omitempty"` // TODO: Save enum
	CreatedAt time.Time `json:"role,omitempty"` // TODO: Save enum
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
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: hashedPass,
		Role:     Customer,
	}

	return user, nil
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func createHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generating hashed password: %s", err)
	}

	return string(hashedPass), nil
}

package domain

import (
	"fmt"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required,min=1"`
	Password  string    `json:"password,omitempty" validate:"required,min=8"`
	Email     string    `json:"email" validate:"required,email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Role string

const (
	Customer Role = "customer"
	Tenant   Role = "tenant"
	Admin    Role = "admin"
)

func NewUser(name, email, password string) (*User, error) {
	user := &User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     Customer,
	}

	ok, err := user.IsValid()
	if !ok {
		return nil, fmt.Errorf("invalid inputs for user: %s", err)
	}

	hashedPass, err := createHashedPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed creating hashed pass: %s", err)
	}
	user.Password = hashedPass

	return user, nil
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) IsValid() (bool, error) {
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}

	return true, nil
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

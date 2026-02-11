package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUser_IsValid(t *testing.T) {
	user := &User{
		ID:        uuid.UUID{},
		Name:      "Tester",
		Password:  "12345678",
		Email:     "andrewhoj@gmail.com",
		Role:      Customer,
		CreatedAt: time.Time{},
	}

	ok, err := user.IsValid()
	if !ok {
		t.Errorf("didn't expect error while validating user, got err: %s", err)
	}
}

func TestUser_IsValidInvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		test *User
	}{
		{
			name: "Too short password",
			test: &User{
				ID:        uuid.UUID{},
				Name:      "Tester",
				Password:  "123",
				Email:     "andrewhoj@gmail.com",
				Role:      Customer,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "Wrongly formatted email",
			test: &User{
				ID:        uuid.UUID{},
				Name:      "Tester",
				Password:  "12345678",
				Email:     "andrewhojgmail.com",
				Role:      Customer,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "Wrongly formatted email and too short password",
			test: &User{
				ID:        uuid.UUID{},
				Name:      "Tester",
				Password:  "1234",
				Email:     "andr.com",
				Role:      Customer,
				CreatedAt: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := tt.test.IsValid()
			if ok {
				t.Errorf("expected to be not ok, error: %s", err)
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	pass := "123456789"
	hash, err := createHashedPassword(pass)
	if err != nil {
		t.Errorf("didnt expect creation of hashed password to return error: %s", err)
	}

	user := &User{
		ID:        uuid.UUID{},
		Name:      "Tester",
		Password:  hash,
		Email:     "andrewhoj@gmail.com",
		Role:      Customer,
		CreatedAt: time.Time{},
	}

	if !user.ComparePassword(pass) {
		t.Error("expected passwords to not equal each other")
	}
}

func TestUser_ComparePasswordIncorrectMatch(t *testing.T) {
	hash, err := createHashedPassword("12345678")
	if err != nil {
		t.Errorf("didnt expect creation of hashed password to return error: %s", err)
	}

	user := &User{
		ID:        uuid.UUID{},
		Name:      "Tester",
		Password:  hash,
		Email:     "andrewhoj@gmail.com",
		Role:      Customer,
		CreatedAt: time.Time{},
	}

	if user.ComparePassword("87654321") {
		t.Error("expected passwords to not equal each other")
	}
}

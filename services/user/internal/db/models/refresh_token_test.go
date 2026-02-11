package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRefreshToken_IsValid(t *testing.T) {
	refreshToken := &RefreshToken{
		ID:        "",
		UserID:    uuid.UUID{},
		Token:     "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 60),
	}

	if !refreshToken.IsValid() {
		t.Error("expected token to be valid")
	}
}

func TestRefreshToken_IsValidShouldFail(t *testing.T) {
	refreshToken := &RefreshToken{
		ID:        "",
		UserID:    uuid.UUID{},
		Token:     "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(-time.Hour * 60),
	}

	if refreshToken.IsValid() {
		t.Error("expected token to be valid")
	}
}

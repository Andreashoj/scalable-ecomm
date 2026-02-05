package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string
	UserID    uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewRefreshToken(userID uuid.UUID) *RefreshToken {
	return &RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 720), // 30 days
	}
}

func (r *RefreshToken) ToString() (string, error) {
	claims := struct {
		userID string
		jwt.RegisteredClaims
	}{
		userID: r.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(r.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(r.CreatedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("failed signing access token: %s", err)
	}

	return tokenString, nil
}

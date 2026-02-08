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
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewRefreshToken(userID uuid.UUID) (*RefreshToken, error) {
	refreshToken := &RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 720),
	}

	token, err := CreateToken(userID.String())
	if err != nil {
		return nil, err
	}

	refreshToken.Token = token
	return refreshToken, nil
}

func CreateToken(userID string) (string, error) {
	claims := struct {
		userID string
		jwt.RegisteredClaims
	}{
		userID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("failed signing access token: %s", err)
	}

	return tokenString, nil
}

func (r *RefreshToken) IsValid() bool {
	return time.Now().Before(r.ExpiresAt)
}

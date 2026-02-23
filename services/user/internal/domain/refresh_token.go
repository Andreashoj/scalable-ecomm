package domain

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string    `json:"id,omitempty"`
	UserID    uuid.UUID `json:"userId,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewRefreshToken(userID uuid.UUID) (*RefreshToken, error) {
	refreshToken := &RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 720),
	}

	token, err := CreateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken.Token = token
	return refreshToken, nil
}

func CreateToken(rToken *RefreshToken) (string, error) {
	claims := struct {
		userID string
		jwt.RegisteredClaims
	}{
		userID: rToken.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(rToken.CreatedAt),
			ExpiresAt: jwt.NewNumericDate(rToken.ExpiresAt),
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

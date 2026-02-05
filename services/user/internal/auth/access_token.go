package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	UserID string
	jwt.RegisteredClaims
}

func createAccessToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 15)
	claims := AccessToken{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte("secret")) // TODO: Update with actual secret
	if err != nil {
		return "", fmt.Errorf("failed signing access token: %s", err)
	}

	return tokenString, nil
}

func parseAccessToken(token string) (*AccessToken, error) {
	var accessToken AccessToken
	err := json.Unmarshal([]byte(token), &accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed parsing access token: %s", err)
	}

	return &accessToken, nil
}

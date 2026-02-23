package pgk

import (
	"andreasho/scalable-ecomm/pgk/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	UserID string `json:"userId,omitempty"`
	Role   string `json:"role,omitempty"`
}

func CreateAccessToken(userID uuid.UUID, role string) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 15)
	claims := jwt.MapClaims{
		"iss": "issuer",
		"exp": expirationTime.Unix(),
		"data": map[string]string{
			"userId": userID.String(),
			"role":   role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret")) // TODO: Update with actual secret
	if err != nil {
		return "", fmt.Errorf("failed signing access token: %s", err)
	}

	return tokenString, nil
}

func parseAccessToken(hashedAccessToken string) (*AccessToken, error) {
	token, err := jwt.Parse(hashedAccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("failed parsing access token: %s", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("couldn't retrieve data from token: %s")
	}

	userID, ok := data["userId"].(string)
	if !ok {
		return nil, errors.New("user id not found or wrong type")
	}

	role, ok := data["role"].(string)
	if !ok {
		return nil, errors.New("role not found or wrong type")
	}

	return &AccessToken{UserID: userID, Role: role}, nil
}

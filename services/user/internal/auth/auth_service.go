package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(payload dto.RegisterUserDTO) error
	Login(payload dto.LoginRequestDTO) (string, string, error)
}

type authService struct {
	userRepo         repos.UserRepo
	refreshTokenRepo repos.RefreshTokenRepo
	logger           pgk.Logger
}

func (a *authService) RegisterUser(payload dto.RegisterUserDTO) error {
	user, err := models.NewUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		return fmt.Errorf("failed creating user: %s", err)
	}
	err = a.userRepo.Save(user)
	if err != nil {
		return fmt.Errorf("failed saving user: %s", err)
	}

	a.logger.Info("Created user with ID: %s", user.GetID())
	return nil
}

func (a *authService) Login(payload dto.LoginRequestDTO) (string, string, error) {
	user, err := a.userRepo.Find(payload.Email)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find user with email: %s", payload.Email)
	}

	validLogin := user.ComparePassword(payload.Password)
	if validLogin {
		refreshToken := models.NewRefreshToken(user.GetID())
		a.refreshTokenRepo.Save(refreshToken)
		refreshTokenWithClaims, err := refreshToken.ToString()
		if err != nil {
			return "", "", fmt.Errorf("failed creating claims for refresh token: %s", err)
		}

		accessTokenWithClaims, err := createAccessToken(user.GetID())
		if err != nil {
			return "", "", fmt.Errorf("failed creating access token: %s", err)
		}

		return refreshTokenWithClaims, accessTokenWithClaims, nil
	}

	return "", "", nil
}

func createAccessToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 15)
	claims := struct {
		userID string
		jwt.RegisteredClaims
	}{
		userID: userID.String(),
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

func NewAuthService(logger pgk.Logger, userRepo repos.UserRepo, token repos.RefreshTokenRepo) AuthService {
	return &authService{
		logger:           logger,
		refreshTokenRepo: token,
		userRepo:         userRepo,
	}
}

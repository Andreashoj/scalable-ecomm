package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	models2 "andreasho/scalable-ecomm/services/user/internal/domain"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"errors"
	"fmt"
)

type AuthService interface {
	RegisterUser(payload dto.RegisterUserDTO) error
	Login(payload dto.LoginRequestDTO) (string, string, error)
	GetUser(userID string) (*models2.User, error)
	InvalidateRefreshToken(refreshToken string) error
	RefreshAccessToken(token string) (string, error)
}

type authService struct {
	userRepo         repos.UserRepo
	refreshTokenRepo repos.RefreshTokenRepo
	logger           pgk.Logger
}

func (a *authService) RefreshAccessToken(refreshToken string) (string, error) {
	rToken, err := a.refreshTokenRepo.Find(refreshToken)
	if err != nil {
		return "", err
	}

	if !rToken.IsValid() {
		return "", errors.New("refresh token expired")
	}

	accessToken, err := createAccessToken(rToken.UserID)
	if err != nil {
		return "", fmt.Errorf("failed creating access token: %s", err)
	}

	return accessToken, nil
}

func (a *authService) GetUser(userID string) (*models2.User, error) {
	user, err := a.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed finding user with error: %s", err)
	}

	return user, nil
}

func (a *authService) RegisterUser(payload dto.RegisterUserDTO) error {
	user, err := models2.NewUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		return fmt.Errorf("failed creating user: %s", err)
	}

	err = a.userRepo.Save(user)
	if err != nil {
		return fmt.Errorf("failed saving user: %s", err)
	}

	a.logger.Info("Created user with ID", "info", user.GetID())
	return nil
}

func (a *authService) Login(payload dto.LoginRequestDTO) (string, string, error) {
	user, err := a.userRepo.FindByEmail(payload.Email)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find user with email: %s", payload.Email)
	}

	validLogin := user.ComparePassword(payload.Password)
	if validLogin {
		refreshToken, err := models2.NewRefreshToken(user.GetID())
		if err != nil {
			return "", "", fmt.Errorf("failed creating refresh token: %s", err)
		}
		err = a.refreshTokenRepo.Save(refreshToken)
		if err != nil {
			return "", "", fmt.Errorf("failed saving refresh token: %s", err)
		}

		accessTokenWithClaims, err := createAccessToken(user.GetID())
		if err != nil {
			return "", "", fmt.Errorf("failed creating access token: %s", err)
		}

		return refreshToken.Token, accessTokenWithClaims, nil
	}

	return "", "", errors.New("password did not match")
}

func (a *authService) InvalidateRefreshToken(refreshToken string) error {
	err := a.refreshTokenRepo.Delete(refreshToken)
	if err != nil {
		return fmt.Errorf("failed deleting refresh token: %s", err)
	}

	return nil
}

func NewAuthService(logger pgk.Logger, userRepo repos.UserRepo, token repos.RefreshTokenRepo) AuthService {
	return &authService{
		logger:           logger,
		refreshTokenRepo: token,
		userRepo:         userRepo,
	}
}

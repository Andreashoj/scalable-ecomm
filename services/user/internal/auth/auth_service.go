package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"fmt"
)

type AuthService interface {
	RegisterUser(payload dto.RegisterUserDTO) error
	Login(payload dto.LoginRequestDTO) (string, string, error)
	Logout(payload dto.LoginRequestDTO) (string, string, error)
	GetUser(userID string) (*models.User, error)
}

type authService struct {
	userRepo         repos.UserRepo
	refreshTokenRepo repos.RefreshTokenRepo
	logger           pgk.Logger
}

func (a *authService) GetUser(userID string) (*models.User, error) {
	user, err := a.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed finding user with error: %s", err)
	}

	return user, nil
}

func (a *authService) Logout(payload dto.LoginRequestDTO) (string, string, error) {
	//TODO implement me
	panic("implement me")
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
	user, err := a.userRepo.FindByEmail(payload.Email)
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

func NewAuthService(logger pgk.Logger, userRepo repos.UserRepo, token repos.RefreshTokenRepo) AuthService {
	return &authService{
		logger:           logger,
		refreshTokenRepo: token,
		userRepo:         userRepo,
	}
}

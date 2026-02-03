package auth

import (
	"andreasho/scalable-ecomm/db/models"
	"andreasho/scalable-ecomm/db/repos"
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"fmt"
)

type AuthService interface {
	RegisterUser(payload dto.RegisterUserDTO) (string, error)
	Login(payload dto.LoginRequestDTO) (*models.User, string, error)
}

type authService struct {
	userRepo        repos.UserRepo
	accessTokenRepo repos.AccessTokenRepo
	refreshTokenRepo repos.RefreshTokenRepo
	logger          pgk.Logger
}

func (a *authService) RegisterUser(payload dto.RegisterUserDTO) (string, error) {
	user, err := models.NewUser(payload.Name, payload.Email, payload.Password)
	if err != nil {
		return "", fmt.Errorf("failed creating user: %s", err)
	}

	_ := user.CreateRefreshToken()
	a.userRepo.Save(user)

	accessToken, accessTokenID := models.NewAccessToken()
	a.accessTokenRepo.Save(accessToken)

	a.logger.Info("Created user with ID: %v", user.GetID())
	return accessTokenID, nil
}

func (a *authService) Login(payload dto.LoginRequestDTO) (*models.User, string, error) {
	// Find user based off email
	user, err := a.userRepo.Find(payload.Email)
	if err != nil {
		return nil, "", fmt.Errorf("couldn't find user with email: %s", payload.Email)
	}

	validLogin := user.ComparePassword(payload.Password)
	if validLogin {
		// Save accesstoken/refresh token
		refreshToken := models.NewRefreshToken(user.GetID())
		a.refreshTokenRepo.Save(refreshToken)

		accessToken := models.NewAccessToken(re)
		return user,
	}

	return nil, "", nil
}

func NewAuthService(logger pgk.Logger, userRepo repos.UserRepo, accessTokenRepo repos.AccessTokenRepo) AuthService {
	return &authService{
		logger:          logger,
		userRepo:        userRepo,
		accessTokenRepo: accessTokenRepo,
	}
}

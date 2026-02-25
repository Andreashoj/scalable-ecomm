package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/domain"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"errors"
	"fmt"
	"time"
)

type AuthService interface {
	UpdateUser(payload dto.UpdateUserRequest) (*domain.User, error)
	RegisterUser(payload dto.RegisterUser) error
	Login(payload dto.LoginRequest) (string, string, error)
	GetUser(userID string) (*domain.User, error)
	InvalidateRefreshToken(refreshToken string) error
	RefreshAccessToken(token string) (string, error)
}

type authService struct {
	userRepo         repos.UserRepo
	refreshTokenRepo repos.RefreshTokenRepo
	logger           pgk.Logger
}

func (a *authService) UpdateUser(payload dto.UpdateUserRequest) (*domain.User, error) {
	user, err := a.GetUser(payload.UserID.String())
	if err != nil {
		return nil, err
	}

	hashedPass, err := domain.CreateHashedPassword(payload.Password)
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = hashedPass
	user.Role = payload.Role

	ok, err := user.IsValid()
	if !ok {
		return nil, fmt.Errorf("invalid user inputs while updating: %v", err)
	}

	err = a.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authService) RefreshAccessToken(refreshToken string) (string, error) {
	rToken, err := a.refreshTokenRepo.Find(refreshToken)
	if err != nil {
		return "", err
	}

	if !rToken.IsValid() {
		return "", errors.New("refresh token expired")
	}

	user, err := a.userRepo.FindByID(rToken.UserID.String())
	if err != nil {
		return "", fmt.Errorf("failed finding user by id: %v", err)
	}

	accessToken, err := pgk.CreateAccessToken(rToken.UserID, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("failed creating access token: %s", err)
	}

	return accessToken, nil
}

func (a *authService) GetUser(userID string) (*domain.User, error) {
	user, err := a.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed finding user with error: %s", err)
	}

	return user, nil
}

func (a *authService) RegisterUser(payload dto.RegisterUser) error {
	user, err := domain.NewUser(payload.Name, payload.Email, payload.Password)
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

func (a *authService) Login(payload dto.LoginRequest) (string, string, error) {
	user, err := a.userRepo.FindByEmail(payload.Email)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find user with email: %s", payload.Email)
	}

	validLogin := user.ComparePassword(payload.Password)
	if validLogin {
		refreshToken, err := domain.NewRefreshToken(user.GetID())
		if err != nil {
			return "", "", fmt.Errorf("failed creating refresh token: %s", err)
		}
		err = a.refreshTokenRepo.Save(refreshToken)
		if err != nil {
			return "", "", fmt.Errorf("failed saving refresh token: %s", err)
		}

		accessTokenWithClaims, err := pgk.CreateAccessToken(user.GetID(), string(user.Role))
		if err != nil {
			return "", "", fmt.Errorf("failed creating access token: %s", err)
		}

		return refreshToken.Token, accessTokenWithClaims, nil
	}

	return "", "", errors.New("password did not match")
}

func (a *authService) InvalidateRefreshToken(refreshToken string) error {
	token, err := a.refreshTokenRepo.Find(refreshToken)
	if err != nil {
		return fmt.Errorf("failed finding refresh token from payload: %v", err)
	}

	if !token.ExpiresAt.After(time.Now()) {
		return errors.New("token is expired")
	}

	err = a.refreshTokenRepo.Delete(refreshToken)
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

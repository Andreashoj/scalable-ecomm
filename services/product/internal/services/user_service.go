package services

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type UserService interface {
	IsAdmin(authorizationHeader string) (bool, error)
}

type userService struct {
	baseURL    string
	httpClient *http.Client
	logger     pgk.Logger
}

func (u *userService) IsAdmin(authorizationHeader string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/me", u.baseURL), nil)
	if err != nil {
		return false, fmt.Errorf("failed validation request: %v", err)
	}

	req.Header.Set("Authorization", authorizationHeader)

	res, err := u.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed making authorization validation request: %v", err)
	}

	var user dto.UserValidationRequest
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return false, fmt.Errorf("failed decoding user validation response: %v", err)
	}

	if user.Role == "admin" {
		return true, nil
	}

	return false, nil
}

func NewUserService() UserService {
	return &userService{
		baseURL:    os.Getenv("USER_SERVICE_URL"),
		httpClient: &http.Client{Timeout: time.Minute * 60},
	}
}

package dto

import (
	"andreasho/scalable-ecomm/services/user/internal/domain"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RegisterUser struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

type UpdateUserRequest struct {
	UserID   uuid.UUID   `json:"userId,omitempty"`
	Name     string      `json:"name,omitempty"`
	Email    string      `json:"email,omitempty"`
	Password string      `json:"password,omitempty"`
	Role     domain.Role `json:"role,omitempty"`
}

type UpdateUserResponse struct {
	Name  string      `json:"name,omitempty"`
	Email string      `json:"email,omitempty"`
	Role  domain.Role `json:"role,omitempty"`
}

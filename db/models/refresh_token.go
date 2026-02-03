package models

import "github.com/google/uuid"

type RefreshToken struct {
	ID     string
	userID uuid.UUID
}

func NewRefreshToken(userID uuid.UUID) *RefreshToken {
	return &RefreshToken{}
}

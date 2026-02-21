package dto

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Name       string      `json:"name"`
	ProductIDs []uuid.UUID `json:"products"`
}

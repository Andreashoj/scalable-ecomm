package dto

import (
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name       string      `json:"name"`
	Price      float64     `json:"price"`
	Categories []uuid.UUID `json:"categories"`
}

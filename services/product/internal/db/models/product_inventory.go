package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductInventory struct {
	ID        uuid.UUID `json:"ID,omitempty"`
	Amount    int       `json:"amount,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	ProductID uuid.UUID `json:"productID,omitempty"`
}

func NewProductInventory() *ProductInventory {
	return &ProductInventory{}
}

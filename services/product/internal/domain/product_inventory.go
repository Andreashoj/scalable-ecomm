package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductInventory struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Amount    int       `json:"amount,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	ProductID uuid.UUID `json:"productId,omitempty"`
}

func NewProductInventory() *ProductInventory {
	return &ProductInventory{}
}

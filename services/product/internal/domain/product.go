package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Price      float64    `json:"price,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
	}
}

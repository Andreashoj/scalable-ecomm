package models

import "github.com/google/uuid"

type Product struct {
	ID         uuid.UUID  `json:"ID,omitempty"`
	Name       string     `json:"name,omitempty"`
	Price      int        `json:"price,omitempty"`
	Categories []Category `json:"categories,omitempty"`
}

func NewProduct() *Product {
	return &Product{}
}

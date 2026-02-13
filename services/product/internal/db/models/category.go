package models

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `json:"ID,omitempty"`
	Name     string    `json:"name,omitempty"`
	Products []Product `json:"products,omitempty"`
}

func NewCategory() *Category {
	return &Category{}
}

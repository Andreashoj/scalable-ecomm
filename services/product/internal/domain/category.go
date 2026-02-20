package domain

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Products []Product `json:"products,omitempty"`
}

func NewCategory(name string) *Category {
	return &Category{
		ID:   uuid.New(),
		Name: name,
	}
}

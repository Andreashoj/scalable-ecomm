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

func (c *Category) AddProduct(product *Product) {
	c.Products = append(c.Products, *product)
}

func (c *Category) AddProducts(products []Product) {
	for _, product := range products {
		c.Products = append(c.Products, product)
	}
}

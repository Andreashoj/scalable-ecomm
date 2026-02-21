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

type ProductCategory struct {
	ProductID  string `db:"product_id"`
	CategoryID string `db:"category_id"`
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
	}
}

func (p *Product) AddCategory(category *Category) {
	p.Categories = append(p.Categories, *category)
}

func (p *Product) RemoveCategory(category *Category) {
	for i, c := range p.Categories {
		if c.ID == category.ID {
			p.Categories = append(p.Categories[:i], p.Categories[i+1:]...)
			break
		}
	}
}

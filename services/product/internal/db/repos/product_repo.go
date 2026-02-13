package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"database/sql"
	"fmt"
)

type ProductRepo interface {
	GetProducts() ([]models.Product, error)
}

type productRepo struct {
	DB *sql.DB
}

func (p *productRepo) GetProducts() ([]models.Product, error) {
	rows, err := p.DB.Query(`SELECT id, name, price FROM product`)

	if err != nil {
		return nil, fmt.Errorf("failed querying db for products: %s", err)
	}

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, fmt.Errorf("failed scanning product: %s", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func NewProductRepo(DB *sql.DB) ProductRepo {
	return &productRepo{
		DB: DB,
	}
}

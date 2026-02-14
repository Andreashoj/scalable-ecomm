package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ProductRepo interface {
	GetProducts() ([]models.Product, error)
	Find(id uuid.UUID) (*models.Product, error)
	Save(product *models.Product) error
}

type productRepo struct {
	DB *sql.DB
}

func (p *productRepo) Find(id uuid.UUID) (*models.Product, error) {
	rows, err := p.DB.Query(`
		SELECT p.id, p.name, p.price, c.id, c.name
		FROM product p
		LEFT JOIN product_category
		ON p.id= product_category.product_id
		LEFT JOIN category c
		ON product_category.category_id = c.id
		WHERE p.id = $1;`, id.String())
	if err != nil {
		return nil, fmt.Errorf("failed querying product with id: %s got err %s", id.String(), err)
	}

	var product models.Product
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("failed mapping product query: %s", err)
		}

		product.Categories = append(product.Categories, category)
	}

	return &product, nil
}

func (p *productRepo) Save(product *models.Product) error {
	//TODO implement me
	panic("implement me")
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

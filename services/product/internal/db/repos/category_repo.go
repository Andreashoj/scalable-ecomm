package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryRepo interface {
	GetAll() ([]domain.Category, error)
	Save(category *domain.Category, products []uuid.UUID) error
}

type categoryRepo struct {
	DB *sqlx.DB
}

func (c *categoryRepo) GetAll() ([]domain.Category, error) {
	rows, err := c.DB.Query(`SELECT id, name FROM category`)
	if err != nil {
		return nil, fmt.Errorf("failed querying categories: %v", err)
	}

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("failed mapping category row: %v", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (c *categoryRepo) Save(category *domain.Category, products []uuid.UUID) error {
	tx, err := c.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed starting transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO category (id, name) VALUES ($1, $2)`, category.ID, category.Name)
	if err != nil {
		return fmt.Errorf("failed inserting category into db: %v", err)
	}

	var productCategories []domain.ProductCategory
	for _, product := range products {
		productCategories = append(productCategories, domain.ProductCategory{
			ProductID:  product.String(),
			CategoryID: category.ID.String(),
		})
	}

	for _, pc := range productCategories {
		_, err = tx.NamedExec(`INSERT INTO product_category (product_id, category_id) VALUES (:product_id, :category_id)`, pc)
		if err != nil {
			return fmt.Errorf("failed inserting product_category record: %v", err)
		}
	}

	return tx.Commit()
}

func NewCategoryRepo(DB *sqlx.DB) CategoryRepo {
	return &categoryRepo{
		DB: DB,
	}
}

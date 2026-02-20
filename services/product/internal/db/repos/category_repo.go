package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoryRepo interface {
	GetAll() ([]domain.Category, error)
	Save(category *domain.Category) error
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

func (c *categoryRepo) Save(category *domain.Category) error {
	_, err := c.DB.Exec(`INSERT INTO category (id, name) VALUES ($1, $2)`, category.ID, category.Name)
	if err != nil {
		return fmt.Errorf("failed inserting category into db: %v", err)
	}

	return nil
}

func NewCategoryRepo(DB *sqlx.DB) CategoryRepo {
	return &categoryRepo{
		DB: DB,
	}
}

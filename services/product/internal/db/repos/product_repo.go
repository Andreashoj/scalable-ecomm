package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	GetProducts(search *domain.ProductSearch) ([]domain.Product, error)
	GetProductsByCategory(categoryID uuid.UUID) ([]domain.Product, error)
	Find(id uuid.UUID) (*domain.Product, error)
	Save(product *domain.Product, categories []uuid.UUID) error
}

type productRepo struct {
	DB *sqlx.DB
}

func (p *productRepo) GetProductsByCategory(categoryID uuid.UUID) ([]domain.Product, error) {
	rows, err := p.DB.Query(`SELECT p.id, p.name, p.price FROM product_category pc JOIN product p on p.id = pc.product_id WHERE category_id = $1`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed querying products: %v", err)
	}

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed mapping products: %v", err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *productRepo) Find(id uuid.UUID) (*domain.Product, error) {
	rows, err := p.DB.Query(`
		SELECT p.id, p.name, p.price, p.created_at, c.id, c.name
		FROM product p
		LEFT JOIN product_category
		ON p.id= product_category.product_id
		LEFT JOIN category c
		ON product_category.category_id = c.id
		WHERE p.id = $1;`, id.String())

	if err != nil {
		return nil, fmt.Errorf("failed querying product with id: %s got err %w", id.String(), err)
	}

	var product domain.Product
	for rows.Next() {
		var categoryID sql.NullString
		var categoryName sql.NullString
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt, &categoryID, &categoryName)
		if err != nil {
			return nil, fmt.Errorf("failed mapping product query: %w", err)
		}

		if categoryID.Valid {
			categoryUUID, _ := uuid.Parse(categoryID.String)
			category := domain.Category{
				ID:   categoryUUID,
				Name: categoryName.String,
			}
			product.Categories = append(product.Categories, category)
		}
	}

	return &product, nil
}

func (p *productRepo) Save(product *domain.Product, categories []uuid.UUID) error {
	tx, err := p.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed starting transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO product (id, name, price, created_at) VALUES ($1, $2, $3, $4)`, product.ID, product.Name, product.Price, product.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed creating product: %w", err)
	}

	var productCategory []domain.ProductCategory
	for _, category := range categories {
		productCategory = append(productCategory, domain.ProductCategory{ProductID: product.ID.String(), CategoryID: category.String()})
	}

	for _, pc := range productCategory {
		_, err = tx.NamedExec(`INSERT INTO product_category (product_id, category_id) VALUES (:product_id, :category_id)`, pc)
		if err != nil {
			return fmt.Errorf("failed inserting product_category record: %v", err)
		}
	}

	return tx.Commit()
}

func (p *productRepo) GetProducts(productSearch *domain.ProductSearch) ([]domain.Product, error) {
	rows, err := p.DB.Query(
		fmt.Sprintf(`SELECT id, name, price, created_at FROM product ORDER BY %s %s`,
			productSearch.Sort.ToSQL(),
			productSearch.Order.ToSQL()))

	if err != nil {
		return nil, fmt.Errorf("failed querying db for products: %w", err)
	}

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed scanning product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func NewProductRepo(DB *sqlx.DB) ProductRepo {
	return &productRepo{
		DB: DB,
	}
}

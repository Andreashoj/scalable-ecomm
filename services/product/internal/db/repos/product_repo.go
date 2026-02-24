package repos

import (
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
	rows, err := p.DB.Query(`
		SELECT p.id, p.name, p.price
		FROM product_category pc 
		JOIN product p 
		on p.id = pc.product_id 
		WHERE category_id = $1`, categoryID)

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
	var rows []productCategoryRow
	var q sq.SelectBuilder
	if len(productSearch.Filter) > 0 {
		q = sq.Select("p.id AS product_id, p.name AS product_name, p.price, c.id AS category_id, c.name AS category_name").
			From("category c").
			LeftJoin("product_category pc ON pc.category_id = c.id").
			LeftJoin("product p ON p.id = pc.product_id").
			Where(sq.Eq{"c.id": productSearch.Filter}).
			OrderBy(fmt.Sprintf("%s %s", productSearch.Sort.ToSQL(), productSearch.Order.ToSQL()))
	} else {
		q = sq.Select("p.id AS product_id, p.name AS product_name, p.price, c.id AS category_id, c.name AS category_name").
			From("product p").
			LeftJoin("product_category pc ON pc.product_id = p.id").
			LeftJoin("category c ON c.id = pc.category_id").
			OrderBy(fmt.Sprintf("%s %s", productSearch.Sort.ToSQL(), productSearch.Order.ToSQL()))
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed converting query into sql: %v", err)
	}

	err = p.DB.Select(&rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed querying db for products: %w", err)
	}

	var productMap = make(map[uuid.UUID]domain.Product)
	for _, row := range rows {
		product, exists := productMap[row.ProductID]
		if !exists {
			productMap[row.ProductID] = domain.Product{
				ID:    row.ProductID,
				Name:  row.ProductName,
				Price: row.ProductPrice,
			}
		}

		var category domain.Category
		if row.CategoryID.Valid {
			category = domain.Category{
				ID:   row.CategoryID.UUID,
				Name: row.CategoryName.String,
			}
		}

		product.Categories = append(product.Categories, category)
		productMap[row.ProductID] = product
	}

	var products []domain.Product
	for _, product := range productMap {
		products = append(products, product)
	}

	return products, nil
}

func NewProductRepo(DB *sqlx.DB) ProductRepo {
	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &productRepo{
		DB: DB,
	}
}

type productCategoryRow struct {
	ProductID    uuid.UUID      `db:"product_id"`
	ProductName  string         `db:"product_name"`
	ProductPrice float64        `db:"price"`
	CategoryID   uuid.NullUUID  `db:"category_id"`
	CategoryName sql.NullString `db:"category_name"`
}

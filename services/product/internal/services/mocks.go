package services

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"fmt"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func SetupProductCatalogService(t *testing.T, productsToAdd int) (ProductCatalogService, repos.ProductRepo, repos.CategoryRepo) {
	db := pgk.SetupTestDB(t, "../services/product/internal/db/migrations")
	DB := sqlx.NewDb(db, "postgres")
	prodRepo := repos.NewProductRepo(DB)
	categoryRepo := repos.NewCategoryRepo(DB)

	for i := 0; i < productsToAdd; i++ {
		product := &domain.Product{
			ID:         uuid.New(),
			Name:       fmt.Sprintf("product-%v", i),
			Price:      float64(i) * 10,
			Categories: nil,
			CreatedAt:  time.Now().Add(-time.Minute * time.Duration(i)),
		}
		err := prodRepo.Save(product, nil)
		if err != nil {
			t.Fatalf("failed creating product: :%v", err)
		}
	}

	return &productCatalogService{productRepo: prodRepo, categoryRepo: categoryRepo}, prodRepo, categoryRepo
}

type MockUserService struct {
	Admin bool
}

func (u *MockUserService) IsAdmin(accessToken string) (bool, error) {
	if u.Admin {
		return true, nil
	}
	return false, nil
}

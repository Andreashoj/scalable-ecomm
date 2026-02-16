package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"fmt"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

func SetupProductCatalogService(t *testing.T, productsToAdd int) (ProductCatalogService, repos.ProductRepo) {
	DB := db.SetupTestDB(t)
	prodRepo := repos.NewProductRepo(DB)

	for i := 0; i < productsToAdd; i++ {
		product := &domain.Product{
			ID:         uuid.New(),
			Name:       fmt.Sprintf("product-%v", i),
			Price:      float64(i) * 10,
			Categories: nil,
			CreatedAt:  time.Now().Add(-time.Minute * time.Duration(i)),
		}
		err := prodRepo.Save(product)
		fmt.Println("no fail?")
		if err != nil {
			t.Fatalf("failed creating product: :%v", err)
		}
	}

	return &productCatalogService{productRepo: prodRepo}, prodRepo
}

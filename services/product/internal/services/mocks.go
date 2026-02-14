package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type productRepoMock struct {
	products map[string]models.Product
}

func (p *productRepoMock) GetProducts() ([]models.Product, error) {
	var products []models.Product
	for _, product := range p.products {
		products = append(products, product)
	}

	return products, nil
}

func (p *productRepoMock) Save(product *models.Product) error {
	p.products[product.ID.String()] = *product
	return nil
}

func (p *productRepoMock) Find(id uuid.UUID) (*models.Product, error) {
	product, exists := p.products[id.String()]
	if !exists {
		return nil, errors.New("couldn't find product with that id")
	}

	return &product, nil
}

func SetupProductCatalogService() (ProductCatalogService, repos.ProductRepo) {
	products := make(map[string]models.Product)
	for i := 0; i < 5; i++ {
		product := models.NewProduct(fmt.Sprintf("product-%v", i), float64(10*i))
		products[product.ID.String()] = *product
	}

	productRepo := &productRepoMock{
		products: products,
	}

	return &productCatalogService{productRepo: productRepo}, productRepo
}

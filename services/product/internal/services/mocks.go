package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"fmt"
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

func SetupProductCatalogService() ProductCatalogService {
	products := make(map[string]models.Product)
	for i := 0; i < 5; i++ {
		product := models.NewProduct(fmt.Sprintf("product-%v", i), float64(10*i))
		products[product.ID.String()] = *product
	}

	return &productCatalogService{productRepo: &productRepoMock{
		products: products,
	}}
}

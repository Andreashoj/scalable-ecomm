package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"fmt"

	"github.com/google/uuid"
)

type ProductCatalogService interface {
	GetProducts() ([]models.Product, error)
	GetProduct(id uuid.UUID) (*models.Product, error)
}

type productCatalogService struct {
	productRepo repos.ProductRepo
}

func (p *productCatalogService) GetProduct(id uuid.UUID) (*models.Product, error) {
	product, err := p.productRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving product: %s", err)
	}

	return product, nil
}

func (p *productCatalogService) GetProducts() ([]models.Product, error) {
	products, err := p.productRepo.GetProducts()
	if err != nil {
		return nil, fmt.Errorf("failed getting products from db: %s", err)
	}

	return products, nil
}

func NewProductCatalogService(productRepo repos.ProductRepo) ProductCatalogService {
	return &productCatalogService{
		productRepo: productRepo,
	}
}

package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

type ProductCatalogService interface {
	GetProducts(search *domain.ProductSearch) ([]domain.Product, error)
	GetProduct(id uuid.UUID) (*domain.Product, error)
}

type productCatalogService struct {
	productRepo repos.ProductRepo
}

func (p *productCatalogService) GetProduct(id uuid.UUID) (*domain.Product, error) {
	product, err := p.productRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving product: %s", err)
	}

	return product, nil
}

func (p *productCatalogService) GetProducts(productSearch *domain.ProductSearch) ([]domain.Product, error) {
	if !productSearch.Order.IsValid() {
		productSearch.Order = domain.OrderAscending
	}

	if !productSearch.Sort.IsValid() {
		productSearch.Sort = domain.SortDate
	}

	products, err := p.productRepo.GetProducts(productSearch)
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

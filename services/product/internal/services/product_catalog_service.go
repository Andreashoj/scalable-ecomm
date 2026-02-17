package services

import (
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"andreasho/scalable-ecomm/services/product/internal/dto"
	"fmt"

	"github.com/google/uuid"
)

type ProductCatalogService interface {
	GetProducts(search *domain.ProductSearch) ([]domain.Product, error)
	GetProduct(id uuid.UUID) (*domain.Product, error)
	CreateProduct(payload dto.CreateProductRequest) (*domain.Product, error)
}

type productCatalogService struct {
	productRepo repos.ProductRepo
}

func (p *productCatalogService) CreateProduct(payload dto.CreateProductRequest) (*domain.Product, error) {
	product := domain.NewProduct(payload.Name, payload.Price)
	err := p.productRepo.Save(product)

	if err != nil {
		return nil, fmt.Errorf("failed saving product: %v", err)
	}

	return product, nil
}

func (p *productCatalogService) GetProduct(id uuid.UUID) (*domain.Product, error) {
	product, err := p.productRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving product: %v", err)
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
		return nil, fmt.Errorf("failed getting products from db: %v", err)
	}

	return products, nil
}

func NewProductCatalogService(productRepo repos.ProductRepo) ProductCatalogService {
	return &productCatalogService{
		productRepo: productRepo,
	}
}

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
	CreateCategory(payload dto.CreateCategoryRequest) (*domain.Category, error)
	GetCategories() ([]domain.Category, error)
}

type productCatalogService struct {
	productRepo  repos.ProductRepo
	categoryRepo repos.CategoryRepo
}

func (p *productCatalogService) GetCategories() ([]domain.Category, error) {
	categories, err := p.categoryRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed getting categories: %v", err)
	}

	return categories, nil
}

func (p *productCatalogService) CreateCategory(payload dto.CreateCategoryRequest) (*domain.Category, error) {
	category := domain.NewCategory(payload.Name)
	err := p.categoryRepo.Save(category, nil)
	if err != nil {
		return nil, fmt.Errorf("failed saving category: %v", err)
	}

	return category, nil
}

func (p *productCatalogService) CreateProduct(payload dto.CreateProductRequest) (*domain.Product, error) {
	product := domain.NewProduct(payload.Name, payload.Price)
	categories := payload.Categories
	err := p.productRepo.Save(product, categories)

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

func NewProductCatalogService(productRepo repos.ProductRepo, categoryRepo repos.CategoryRepo) ProductCatalogService {
	return &productCatalogService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

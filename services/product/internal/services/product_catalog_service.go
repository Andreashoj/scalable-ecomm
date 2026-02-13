package services

type ProductCatalogService interface{}

type productCatalogService struct{}

func NewProductCatalogService() ProductCatalogService {
	return &productCatalogService{}
}

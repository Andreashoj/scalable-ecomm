package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type RouterHandler struct {
	logger                pgk.Logger
	productCatalogService services.ProductCatalogService
}

func StartRouterHandlers(r *chi.Mux, logger pgk.Logger, productCatalogService services.ProductCatalogService) error {
	h := &RouterHandler{
		logger:                logger,
		productCatalogService: productCatalogService,
	}

	r.Get("/products", h.GetProducts)
	r.Get("/product/{id}", h.GetProduct)

	return nil
}

func (h *RouterHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	order := chi.URLParam(r, "order")
	sort := chi.URLParam(r, "sort")

	productSearch := domain.NewProductSearch(order, sort)

	products, err := h.productCatalogService.GetProducts(productSearch)
	if err != nil {
		h.logger.Error("failed getting products", "error", err)
		rest.ErrorResponse(w, 500, "failed retrieving products")
		return
	}

	rest.Response(w, products, 200)
}

func (h *RouterHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("non existent id passed in as query parameter for product", "error", err)
		w.WriteHeader(404)
		return
	}

	product, err := h.productCatalogService.GetProduct(productID)
	if err != nil {
		h.logger.Error("couldn't find product with matching id", "error", err)
		rest.ErrorResponse(w, 404, "couldn't find product matching that description")
		return
	}

	rest.Response(w, product, 200)
}

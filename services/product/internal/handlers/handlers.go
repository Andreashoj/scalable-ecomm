package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"net/http"

	"github.com/go-chi/chi"
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
	return nil
}

func (h *RouterHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productCatalogService.GetProducts()
	if err != nil {
		h.logger.Error("failed getting products", "error", err)
		rest.ErrorResponse(w, 500, "failed retrieving products")
		return
	}

	rest.Response(w, products, 200)
}

package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"andreasho/scalable-ecomm/services/product/internal/dto"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	r.Get("/category", h.GetCategories)

	r.Group(func(g chi.Router) {
		g.Use(pgk.IsAdmin)
		g.Post("/product", h.CreateProduct)
		g.Post("/category", h.CreateCategory)
	})

	return nil
}

func (h *RouterHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productCatalogService.GetCategories()
	if err != nil {
		rest.ErrorResponse(w, 500, errors.ErrorMessage(err.Error()))
		return
	}

	rest.Response(w, categories, 200)
}

func (h *RouterHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateCategoryRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("failed decoding payload", "error", err)
		rest.ErrorResponse(w, 500, errors.ErrorMessage(err.Error()))
		return
	}

	category, err := h.productCatalogService.CreateCategory(payload)
	if err != nil {
		h.logger.Error("failed creating product", "error", err)
		rest.ErrorResponse(w, 500, errors.ErrorMessage(err.Error()))
		return
	}

	rest.Response(w, category, 200)
}

func (h *RouterHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed decoding payload request: %s", err))
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	product, err := h.productCatalogService.CreateProduct(payload)
	if err != nil {
		h.logger.Error("failed creating product: %v", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	rest.Response(w, product, 200)
}

func (h *RouterHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	q := r.URL.Query().Get("q")
	var filters []string
	if q != "" {
		filters = strings.Split(q, ",")
	}

	productSearch := domain.NewProductSearch(order, sort, filters)

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

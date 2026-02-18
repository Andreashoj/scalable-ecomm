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

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type RouterHandler struct {
	logger                pgk.Logger
	productCatalogService services.ProductCatalogService
	userService           services.UserService
}

func StartRouterHandlers(r *chi.Mux, logger pgk.Logger, productCatalogService services.ProductCatalogService, userService services.UserService) error {
	h := &RouterHandler{
		logger:                logger,
		productCatalogService: productCatalogService,
		userService:           userService,
	}

	r.Get("/products", h.GetProducts)
	r.Get("/product/{id}", h.GetProduct)
	r.Post("/product", h.CreateProduct)

	return nil
}

func (h *RouterHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	isAdmin, err := h.userService.IsAdmin(authorizationHeader)
	if err != nil {
		h.logger.Error("failed authorization request while trying to create product", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	if !isAdmin {
		rest.ErrorResponse(w, 401, errors.Unauthorized)
		return
	}

	var payload dto.CreateProductRequest
	err = json.NewDecoder(r.Body).Decode(&payload)
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

package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

// What handlers do we need for produtcs
// /product?sort=>price/date
// /product/:id
// /category?sort=>price/date
// /product => POST

func TestHandler_Product(t *testing.T) {
	r, _, productRepo := handlerSetup()

	product := models.NewProduct("tester", float64(8291))
	productRepo.Save(product)

	req := httptest.NewRequest("GET", fmt.Sprintf("/product/%v", product.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %v", w.Code)
	}

	var responseProduct models.Product
	err := json.Unmarshal([]byte(w.Body.String()), &responseProduct)
	if err != nil {
		t.Errorf("expected response to contain product body, instead got error: %s", err)
	}

	if responseProduct.ID != product.ID {
		t.Errorf("expected ID response id to equal %v, instead got %v", product.ID, responseProduct.ID)
	}
}

func TestHandler_ProductNotFound(t *testing.T) {
	r, _, _ := handlerSetup()

	req := httptest.NewRequest("GET", fmt.Sprintf("/product/random-param"), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %v", w.Code)
	}
}

func TestHandler_Products(t *testing.T) {
	r, _, _ := handlerSetup()
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200 but got: %v", w.Code)
	}

	var products []models.Product
	err := json.Unmarshal([]byte(w.Body.String()), &products)
	if err != nil {
		t.Errorf("expected to unmarshal response but failed with error: %s", err)
	}

	if len(products) == 0 {
		t.Error("expected to have received products but got 0")
	}
}

func handlerSetup() (*chi.Mux, services.ProductCatalogService, repos.ProductRepo) {
	logger := pgk.NewLogger()
	service, productRepo := services.SetupProductCatalogService()
	r := chi.NewRouter()
	StartRouterHandlers(r, logger, service)

	return r, service, productRepo
}

package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

// What handlers do we need for produtcs
// /product [X]
// / query params => ?sort=>price/date []
// /product/:id [X]
// /category?sort=>price/date [X]
// /product => POST

func TestHandler_Product(t *testing.T) {
	r, _, productRepo := handlerSetup(t, 0)

	product := domain.NewProduct("tester", float64(8291))
	err := productRepo.Save(product)
	if err != nil {
		t.Fatalf("failed save here: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/product/%v", product.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %v", w.Code)
	}

	var responseProduct domain.Product
	err = json.Unmarshal([]byte(w.Body.String()), &responseProduct)
	if err != nil {
		t.Errorf("expected response to contain product body, instead got error: %s", err)
	}

	if responseProduct.ID != product.ID {
		t.Errorf("expected ID response id to equal %v, instead got %v", product.ID, responseProduct.ID)
	}
}

func TestHandler_ProductNotFound(t *testing.T) {
	r, _, _ := handlerSetup(t, 0)

	req := httptest.NewRequest("GET", fmt.Sprintf("/product/random-param"), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %v", w.Code)
	}
}

func TestHandler_Products(t *testing.T) {
	productAmount := 5
	r, _, _ := handlerSetup(t, productAmount)
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200 but got: %v", w.Code)
	}

	var products []domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &products)
	if err != nil {
		t.Errorf("expected to unmarshal response but failed with error: %s", err)
	}

	if len(products) != productAmount {
		t.Error("expected to have received products but got 0")
	}
}

func TestHandler_ProductsNotFound(t *testing.T) {
	productAmount := 0
	r, _, _ := handlerSetup(t, productAmount)
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200 but got: %v", w.Code)
	}

	var products []domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &products)
	if err != nil {
		t.Errorf("expected to unmarshal response but failed with error: %s", err)
	}

	if len(products) != productAmount {
		t.Error("expected to have received products but got 0")
	}
}

func TestHandler_ProductsDateAscending(t *testing.T) {
	r, _, _ := handlerSetup(t, 50)
	req := httptest.NewRequest("GET", "/products?sort=date&order=ascending", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %v instead", w.Code)
	}

	var response []domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("failed decoding body response with err: %s", err)
	}

	previousProductDate := response[0].CreatedAt
	for i, product := range response {
		if i == 0 {
			continue
		}

		if product.CreatedAt.Before(previousProductDate) {
			t.Errorf("expected products to be in ascending order. Current product: %s, previous product: %s", product.CreatedAt, previousProductDate)
		}

		previousProductDate = product.CreatedAt
	}
}

func TestHandler_ProductsDateDescending(t *testing.T) {
	r, _, _ := handlerSetup(t, 50)
	req := httptest.NewRequest("GET", "/products?sort=date&order=descending", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %v instead", w.Code)
	}

	var response []domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("failed decoding body response with err: %s", err)
	}

	previousProductDate := response[0].CreatedAt
	for i, product := range response {
		if i == 0 {
			continue
		}

		if product.CreatedAt.After(previousProductDate) {
			t.Errorf("expected products date to be before previous products. Current product: %s, previous product: %s", product.CreatedAt, previousProductDate)
		}

		previousProductDate = product.CreatedAt
	}
}

func TestHandler_ProductSearchInvalidInputs(t *testing.T) {
	r, _, _ := handlerSetup(t, 50)
	req := httptest.NewRequest("GET", "/products?sort=invalidsortingoption&order=invalidngorderopion", nil) // ascending is the default option if order option validation fails
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %v instead", w.Code)
	}

	var response []domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("failed decoding body response with err: %s", err)
	}

	previousProductDate := response[0].CreatedAt
	for i, product := range response {
		if i == 0 {
			continue
		}

		if product.CreatedAt.Before(previousProductDate) {
			t.Errorf("expected products date to be before previous products. Current product: %s, previous product: %s", product.CreatedAt, previousProductDate)
		}

		previousProductDate = product.CreatedAt
	}
}

func handlerSetup(t *testing.T, productsToAdd int) (*chi.Mux, services.ProductCatalogService, repos.ProductRepo) {
	logger := pgk.NewLogger()
	service, productRepo := services.SetupProductCatalogService(t, productsToAdd)
	r := chi.NewRouter()
	StartRouterHandlers(r, logger, service)

	return r, service, productRepo
}

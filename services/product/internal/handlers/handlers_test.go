package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/domain"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestHandler_CreateCategory(t *testing.T) {
	r, _, _, categoryRepo := handlerSetup(t, 0, true)

	body := `{"name": "my-category"}`
	req := httptest.NewRequest("POST", "/category", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200 instead got: %v", w.Code)
	}

	categories, err := categoryRepo.GetAll()
	if err != nil {
		t.Errorf("failed getting categories: %v", err)
	}

	if len(categories) != 1 {
		t.Errorf("expected categories length to be 1 instead got %v", len(categories))
	}
}

func TestHandler_CreateProduct(t *testing.T) {
	r, _, productRepo, _ := handlerSetup(t, 0, true)

	body := `{"name": "VHS player", "price": 2819, "categories": ["electronics"]}`
	req := httptest.NewRequest("POST", "/product", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status 200, instead got %v", w.Code)
	}

	var response domain.Product
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("failed decoding json response to product: %v", err)
	}

	_, err = productRepo.Find(response.ID)
	if err != nil {
		t.Errorf("expected to find posted product in DB, but got: %v", err)
	}
}

func TestHandler_CreateProductUnauthorized(t *testing.T) {
	r, _, _, _ := handlerSetup(t, 0, false)

	body := `{"name": "VHS player", "price": 2819}`
	req := httptest.NewRequest("POST", "/product", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected status 401, instead got %v", w.Code)
	}
}

func TestHandler_Product(t *testing.T) {
	r, _, productRepo, _ := handlerSetup(t, 0, false)

	product := domain.NewProduct("tester", float64(8291))
	err := productRepo.Save(product, nil)
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
	r, _, _, _ := handlerSetup(t, 0, false)

	req := httptest.NewRequest("GET", fmt.Sprintf("/product/random-param"), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %v", w.Code)
	}
}

func TestHandler_Products(t *testing.T) {
	productAmount := 5
	r, _, _, _ := handlerSetup(t, productAmount, false)
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
	r, _, _, _ := handlerSetup(t, productAmount, false)
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
	r, _, _, _ := handlerSetup(t, 50, false)
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
	r, _, _, _ := handlerSetup(t, 50, false)
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
	r, _, _, _ := handlerSetup(t, 50, false)
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

func handlerSetup(t *testing.T, productsToAdd int, isAdmin bool) (*chi.Mux, services.ProductCatalogService, repos.ProductRepo, repos.CategoryRepo) {
	logger := pgk.NewLogger()
	productCatalogService, productRepo, categoryRepo := services.SetupProductCatalogService(t, productsToAdd)
	userService := &services.MockUserService{
		Admin: isAdmin,
	}
	r := chi.NewRouter()
	StartRouterHandlers(r, logger, productCatalogService, userService)

	return r, productCatalogService, productRepo, categoryRepo
}

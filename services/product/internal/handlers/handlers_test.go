package handlers

import (
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// What handlers do we need for produtcs
// /product?sort=>price/date
// /product/:id
// /category?sort=>price/date
// /product => POST

func TestHandler_Products(t *testing.T) {
	service := services.NewProductCatalogService()

	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	service.GetProducts(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200 but got: %v", w.Code)
	}

	var products []models.Products
	err := json.Unmarshal([]byte(w.Body.String()), &products)
	if err != nil {
		t.Errorf("expected to unmarshal response but failed with error: %s", err)
	}

	if len(w.Body.String()) == 0 {
		t.Error("expected to have received products but got 0")
	}
}

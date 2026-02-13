package handlers

import (
	"andreasho/scalable-ecomm/services/product/internal/db/models"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

// What handlers do we need for produtcs
// /product?sort=>price/date
// /product/:id
// /category?sort=>price/date
// /product => POST

func TestHandler_Products(t *testing.T) {
	service := services.SetupProductCatalogService()
	r := chi.NewRouter()
	StartRouterHandlers(r, service)

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

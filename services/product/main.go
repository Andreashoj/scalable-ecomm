package main

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/product/internal/db"
	"andreasho/scalable-ecomm/services/product/internal/db/repos"
	"andreasho/scalable-ecomm/services/product/internal/handlers"
	"andreasho/scalable-ecomm/services/product/internal/services"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// Get ENV's in order

	logger := pgk.NewLogger()
	r := chi.NewRouter()

	DB, err := db.StartDB()
	if err != nil {
		logger.Error("failed starting db: %s", err)
		return
	}

	// repos
	productRepo := repos.NewProductRepo(DB)

	// services
	authService := services.NewProductCatalogService(productRepo)

	handlers.StartRouterHandlers(r, logger, authService)

	http.ListenAndServe(":8080", r)
}

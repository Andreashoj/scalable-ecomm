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

	defer DB.Close()

	// repos
	productRepo := repos.NewProductRepo(DB)
	categoryRepo := repos.NewCategoryRepo(DB)

	// services
	authService := services.NewProductCatalogService(productRepo, categoryRepo)
	userService := services.NewUserService()

	handlers.StartRouterHandlers(r, logger, authService, userService)

	http.ListenAndServe(":8080", r)
}

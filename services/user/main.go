package main

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/db"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// User service MVP

	// Register
	// => Validation checks [X]
	// => Duplicate email [X]
	// => Valid password [X]
	// => Valid email [X]
	// Authentication [X]
	// Check if access token works [X]
	// Logout []
	// Refresh []

	DB, err := db.StartDB()
	defer DB.Close()
	if err != nil {
		fmt.Printf("failed creating connection to DB: %s", err)
		return
	}

	r := chi.NewRouter()
	logger := pgk.NewLogger()

	// repos
	userRepo := repos.NewUserRepo(DB)
	refreshTokenRepo := repos.NewRefreshTokenRepo(DB)

	// services
	authService := auth.NewAuthService(logger, userRepo, refreshTokenRepo)

	// handlers
	handlers.StartRouteHandler(r, logger, authService)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("failed while starting http router: %s", err)
	}
}

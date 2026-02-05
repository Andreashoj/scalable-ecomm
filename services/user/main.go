package main

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/handlers"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/db"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// User service MVP

	// Keep is simple for now:
	// Schemas:
	// User: id, name, email, role, refresh_token (hash), refresh_created_at [X]
	// Access Token: token (hash), created_at, expires_at, fk => refresh_token [X]
	// Models => Save migrations ? [X]
	// Endpoints: []
	// Register [X]
	// Login []
	// Logout []
	// Authentication

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

	handlers.StartRouteHandler(r, logger, authService)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("failed while starting http router: %s", err)
	}
}

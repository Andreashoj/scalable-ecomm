package main

import (
	"andreasho/scalable-ecomm/db/repos"
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/handlers"
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

	// Flow => Register => Login => Save Requests (Authentication guard) => Logout
	r := chi.NewRouter()
	logger := pgk.NewLogger()

	// repos
	userRepo := repos.NewUserRepo()
	accessTokenRepo := repos.NewAccessTokenRepo()

	// services
	authService := auth.NewAuthService(logger, userRepo, accessTokenRepo)

	handlers.StartRouteHandler(r, logger, authService)

	http.ListenAndServe(":8080", r)
}

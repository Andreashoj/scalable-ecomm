package main

import (
	"andreasho/scalable-ecomm/pgk"
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
	// Models => Create migrations ? [X]
	// Endpoints: [X]
	// Register [X]
	// Login [X]
	// Logout [X]
	// Authentication

	// Flow => Register => Login => Create Requests (Authentication guard) => Logout
	logger := pgk.NewLogger()
	r := chi.NewRouter()
	handlers.StartRouteHandler(r, logger)

	http.ListenAndServe(":8080", r)
}

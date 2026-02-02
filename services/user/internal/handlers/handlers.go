package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type routeHandler struct {
	logger pgk.Logger
}

func StartRouteHandler(r *chi.Mux, logger pgk.Logger) {
	routerHandler := &routeHandler{
		logger: logger,
	}

	routerHandler.registerRoutes(r)
}

func (r *routeHandler) registerRoutes(router *chi.Mux) error {
	router.Route("/auth", func(a chi.Router) {
		a.Post("/register", registerUser)
		a.Post("/login", login)
		a.Post("/logout", logout)
		a.Get("/me", me)
	})

	return nil
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// Get name, email and password
	var payload dto.RegisterUserDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		rest.ErrorResponse(w, 500, errors.BadRequest)
	}

	// Create user
	// Create access_token

	// Returns:
	// OK, return access_token

	// Client/ssr side:
	// Sets access token in httpOnly cookie (server-side, or just use it somehow), or header
}

func me(w http.ResponseWriter, r *http.Request) {
	// Use access_token to retrieve user
	// Used for validating a user is logged in

	// Returns:
	// Valid access_token => 200
	// Invalid access_token / none => 401
}

func login(w http.ResponseWriter, r *http.Request) {
	// Accepts email and password
	// Create new refresh_token / access_token

	// Returns:
	// User, access_token
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Accepts access_token

	// Invalid access_token => 401
	// Delete refresh_token => 201
}

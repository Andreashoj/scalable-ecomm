package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type routeHandler struct {
	logger pgk.Logger

	authService auth.AuthService
}

func StartRouteHandler(r *chi.Mux, logger pgk.Logger, authService auth.AuthService) {
	routerHandler := &routeHandler{
		logger:      logger,
		authService: authService,
	}

	routerHandler.registerRoutes(r)
}

func (r *routeHandler) registerRoutes(router *chi.Mux) error {
	router.Route("/auth", func(a chi.Router) {
		a.Post("/register", r.registerUser)
		a.Post("/login", r.login)
		a.Post("/logout", logout)
		a.Get("/me", me)
	})

	return nil
}

func (h *routeHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("failed decoding payload with error: %s", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	accessTokenID, err := h.authService.RegisterUser(payload)
	if err != nil {
		h.logger.Error("failed creating user with error: %s", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	rest.Response(w, accessTokenID, 200)
}

func me(w http.ResponseWriter, r *http.Request) {
	// Use access_token to retrieve user
	// Used for validating a user is logged in

	// Returns:
	// Valid access_token => 200
	// Invalid access_token / none => 401
}

func (h *routeHandler) login(w http.ResponseWriter, r *http.Request) {
	// Accepts email and password
	var payload dto.LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("failed decoding login payload with error: %s", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	user, accessTokenID, err := h.authService.Login(payload)
	if err != nil {
		rest.ErrorResponse(w, 401, errors.BadRequest)
		return
	}

	response := map[string]interface{}{
		"user":          user,
		"accessTokenID": accessTokenID,
	}

	rest.Response(w, response, 200)
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Accepts access_token

	// Invalid access_token => 401
	// Delete refresh_token => 201
}

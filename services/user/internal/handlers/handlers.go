package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"encoding/json"
	"fmt"
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

	err := routerHandler.registerRoutes(r)
	if err != nil {
		fmt.Sprintf("route handler failed with error: %s", err)
		return
	}
}

func (h *routeHandler) registerRoutes(router *chi.Mux) error {
	router.Route("/auth", func(a chi.Router) {
		a.Post("/register", h.registerUser)
		a.Post("/login", h.login)
		a.Post("/refresh", h.refresh)

		a.Group(func(g chi.Router) {
			g.Use(auth.AuthMiddleware)
			g.Post("/logout", h.logout)
			g.Get("/me", h.me)
		})
	})

	return nil
}

func (h *routeHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("failed decoding payload with error", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	err = h.authService.RegisterUser(payload)
	if err != nil {
		h.logger.Error("failed creating user with error", "error", err)
		rest.ErrorResponse(w, 500, errors.ErrorMessage(err.Error()))
		return
	}

	rest.Response(w, "", 201) // TODO: empty responses
}

func (h *routeHandler) me(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.AccessToken)
	userID := claims.UserID
	user, err := h.authService.GetUser(userID)

	if err != nil {
		h.logger.Error("failed getting user: %v", err)
		rest.ErrorResponse(w, 500, errors.NotFound)
		return
	}

	rest.Response(w, user, 200)
}

func (h *routeHandler) login(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("Failed decoding login payload", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	refreshToken, accessToken, err := h.authService.Login(payload)
	if err != nil {
		h.logger.Info("Unauthorized attempt to login", "error", err)
		rest.ErrorResponse(w, 401, errors.BadRequest)
		return
	}

	response := dto.LoginResponseDTO{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	rest.Response(w, response, 200)
}

func (h *routeHandler) logout(w http.ResponseWriter, r *http.Request) {
	var payload dto.LogoutRequestDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	err = h.authService.InvalidateRefreshToken(payload.RefreshToken)
	if err != nil {
		h.logger.Error("failed logging user out", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	rest.Response(w, "", 201)
}

func (h *routeHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var payload dto.RefreshRequestDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("failed decoding refresh request", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	accessToken, err := h.authService.RefreshAccessToken(payload.RefreshToken)
	if err != nil {
		h.logger.Error("failed refreshing access token", "error", err)
		rest.ErrorResponse(w, 500, errors.BadRequest)
		return
	}

	response := dto.RefreshResponseDTO{
		AccessToken: accessToken,
	}

	h.logger.Info("access token successfully refreshed")
	rest.Response(w, response, 200)
}

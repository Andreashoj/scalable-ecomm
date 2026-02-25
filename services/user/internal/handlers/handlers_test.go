package handlers

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/auth"
	"andreasho/scalable-ecomm/services/user/internal/domain"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create endpoint to update user, should be behind admin guard
// Seed Admin user

func TestHandler_UpdateUser(t *testing.T) {
	service, userRepo, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	password := "123456789"
	hashedPass, err := domain.CreateHashedPassword(password)
	if err != nil {
		t.Error(err)
	}

	user := &domain.User{
		ID:        uuid.New(),
		Name:      "tester",
		Password:  hashedPass,
		Email:     "tester@gmail.com",
		Role:      domain.Admin,
		CreatedAt: time.Now(),
	}

	err = userRepo.Save(user)
	if err != nil {
		t.Error(err)
	}
	loginRes := authorizeUser(t, r, user.Email, password)

	updatedName, updatedEmail := "updatedUser", "updated@gmail.com"
	payload := fmt.Sprintf(`{"userId": "%s","name": "%s", "email": "%s", "password": "updated-secret", "role": "admin"}`, user.ID, updatedName, updatedEmail)
	req := httptest.NewRequest("PATCH", "/auth/me", strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginRes.AccessToken))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected status code 200, got %v", w.Code)
	}

	var res dto.UpdateUserResponse
	err = json.Unmarshal([]byte(w.Body.String()), &res)
	if err != nil {
		t.Errorf("failed decoding response: %v", err)
	}

	if res.Name != updatedName && res.Email != updatedEmail {
		t.Errorf("expected response to include updated email and name, instead got: %v and %v", res.Name, res.Email)
	}
}

func TestHandler_RefreshAccessToken(t *testing.T) {
	service, userRepo, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	email := "andrewhoj@gmail.com"
	user, _ := domain.NewUser("andrew", email, pass)
	userRepo.Save(user)

	loginRes := authorizeUser(t, r, email, pass)

	payload := fmt.Sprintf(`{"refreshToken": "%s"}`, loginRes.RefreshToken)
	req := httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginRes.AccessToken))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected res code 201, instead got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "accessToken") {
		t.Errorf("expected access token to be present instead got: %s", w.Body.String())
	}
}

func TestHandler_RefreshExpired(t *testing.T) {
	service, userRepo, refreshTokenRepo := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	email := "andrewhoj@gmail.com"
	user, _ := domain.NewUser("andrew", email, pass)
	userRepo.Save(user)

	expRefreshToken := &domain.RefreshToken{
		ID:        "",
		UserID:    user.ID,
		Token:     "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(-time.Hour * 10),
	}

	expToken, err := domain.CreateToken(expRefreshToken)
	if err != nil {
		t.Errorf("didn't expect error on creating expired refresh token, but got: %s", err)
	}
	expRefreshToken.Token = expToken
	refreshTokenRepo.Save(expRefreshToken)

	payload := fmt.Sprintf(`{"refreshToken": "%s"}`, expRefreshToken.Token)
	req := httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(payload))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected res code 500, instead got %v", w.Code)
	}

	if strings.Contains(w.Body.String(), "accessToken") {
		t.Errorf("expected access token not be present instead got: %s", w.Body.String())
	}
}

func TestHandler_RefreshNotPresent(t *testing.T) {
	service, _, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	payload := `{"refreshToken": ""}`
	req := httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(payload))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected res code 500, instead got %v", w.Code)
	}

	if strings.Contains(w.Body.String(), "accessToken") {
		t.Errorf("expected access token not be present instead got: %s", w.Body.String())
	}
}

func TestHandler_Logout(t *testing.T) {
	service, userRepo, refreshTokenRepo := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	email := "andrewhoj@gmail.com"
	user, _ := domain.NewUser("andrew", email, pass)
	userRepo.Save(user)

	loginRes := authorizeUser(t, r, email, pass)

	payload := fmt.Sprintf(`{"refreshToken": "%s"}`, loginRes.RefreshToken)
	req := httptest.NewRequest("POST", "/auth/logout", strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginRes.AccessToken))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected res code 201, instead got %v", w.Code)
	}

	_, err := refreshTokenRepo.Find(loginRes.RefreshToken)
	if err == nil {
		t.Error("expected error as the refresh token should have been deleted")
	}
}

func TestHandler_LogoutInvalidRefreshToken(t *testing.T) {
	service, userRepo, refreshTokenRepo := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	email := "andrewhoj@gmail.com"
	user, _ := domain.NewUser("andrew", email, pass)
	userRepo.Save(user)

	loginRes := authorizeUser(t, r, email, pass)

	payload := `{"refreshToken": ""}`
	req := httptest.NewRequest("POST", "/auth/logout", strings.NewReader(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginRes.AccessToken))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected res code 500, instead got %v", w.Code)
	}

	_, err := refreshTokenRepo.Find(loginRes.RefreshToken)
	if err != nil {
		t.Error("expected error as the refresh token should have been deleted")
	}
}

func TestHandler_Me(t *testing.T) {
	service, userRepo, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	email := "andrewhoj@gmail.com"
	user, _ := domain.NewUser("andrew", email, pass)
	userRepo.Save(user)

	loginRes := authorizeUser(t, r, email, pass)

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, pass)
	req := httptest.NewRequest("GET", "/auth/me", strings.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginRes.AccessToken))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected res code 200, instead got %v", w.Code)
	}
}

func TestHandler_MeUnauthorized(t *testing.T) {
	service, _, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	req := httptest.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("Authorization", "wrongaccesstoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected res code 401, instead got %v", w.Code)
	}
}

func TestHandler_RegisterUser(t *testing.T) {
	service, _, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	body := `{"name": "anz", "email":"andrewhoj@gmail.com","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected status code 201 instead got: %v", w.Code)
	}
}

func TestHandler_RegisterUserBadFormat(t *testing.T) {
	service, _, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	body := `{"name": "anz", "email":"andrewhoj@gmail.com"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected status code 500 instead got: %v", w.Code)
	}
}

func TestHandler_Login(t *testing.T) {
	service, userRepo, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	pass := "12345678"
	user, _ := domain.NewUser("andrew", "andrewhoj@gmail.com", pass)
	userRepo.Save(user)

	body := fmt.Sprintf(`{"email":"andrewhoj@gmail.com","password":"%s"}`, pass)
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected no errors, instead got: %v", w.Code)
	}
}

func TestHandler_LoginIncorrectPassword(t *testing.T) {
	service, userRepo, _ := auth.SetupAuthService(t)
	r := chi.NewRouter()
	StartRouteHandler(r, pgk.NewLogger(), service)

	user, _ := domain.NewUser("andrew", "andrewhoj@gmail.com", "12345678")
	userRepo.Save(user)

	body := `{"email":"andrewhoj@gmail.com","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == 200 {
		t.Errorf("expected 401, instead got: %v", w.Code)
	}
}

func authorizeUser(t *testing.T, router *chi.Mux, email, pass string) *dto.LoginResponse {
	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, pass)
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, instead got: %v", w.Code)
	}

	var res dto.LoginResponse
	err := json.Unmarshal([]byte(w.Body.String()), &res)

	if err != nil {
		t.Errorf("didn't expect the response to fail with: %s", err)
	}

	return &res
}

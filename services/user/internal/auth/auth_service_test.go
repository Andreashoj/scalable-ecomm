package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/dto"
	"fmt"
	"testing"
	"time"
)

func TestAuthService_LoginSuccess(t *testing.T) {
	service, userRepo, _ := setupAuthService()
	// Precondition
	user, _ := models.NewUser("Andrew", "andrewhoj@gmail.com", "123456789")
	userRepo.Save(user)

	payload := dto.LoginRequestDTO{
		Email:    "andrewhoj@gmail.com",
		Password: "123456789",
	}

	refreshToken, accessToken, err := service.Login(payload)
	if err != nil {
		t.Errorf("expected err to be nil instead got: %s", err)
	}

	if len(refreshToken) == 0 || len(accessToken) == 0 {
		t.Errorf("expected refresh and access token to have a length greater than 0, got refresh length: %v and access length: %v", len(refreshToken), len(accessToken))
	}
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	service, userRepo, _ := setupAuthService()

	user, _ := models.NewUser("andrew", "andrewhoj@gmail.com", "123456789")
	userRepo.Save(user)

	payload := dto.LoginRequestDTO{
		Email:    "andrewhoj@gmail.com",
		Password: "wrongpassword",
	}

	refreshToken, accessToken, err := service.Login(payload)

	if len(refreshToken) != 0 || len(accessToken) != 0 {
		t.Errorf("expected refresh token and access token to be empty strings, instead got refresh: %s and access: %s", refreshToken, accessToken)
	}

	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestAuthService_LoginUserNotFound(t *testing.T) {
	service, _, _ := setupAuthService()

	payload := dto.LoginRequestDTO{
		Email:    "andrewhoj@gmail.com",
		Password: "123456789",
	}

	_, _, err := service.Login(payload)

	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestAuthService_RegisterUser(t *testing.T) {
	service, _, _ := setupAuthService()

	payload := dto.RegisterUserDTO{
		Name:     "andreas",
		Email:    "andrewhoj@gmail.com",
		Password: "123456789",
	}

	err := service.RegisterUser(payload)

	if err != nil {
		t.Errorf("expected error to be nil instead got: %s", err)
	}
}

func TestAuthService_RegisterUserInvalidInputs(t *testing.T) {
	service, _, _ := setupAuthService()

	tests := []struct {
		name string
		test dto.RegisterUserDTO
	}{
		{
			name: "missing email",
			test: dto.RegisterUserDTO{
				Name:     "andreas",
				Email:    "",
				Password: "123456789",
			},
		},
		{
			name: "missing password",
			test: dto.RegisterUserDTO{
				Name:     "andreas",
				Email:    "andrewhoj@gmail.com",
				Password: "",
			},
		},
		{
			name: "missing email and password",
			test: dto.RegisterUserDTO{
				Name:     "andreas",
				Email:    "",
				Password: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RegisterUser(tt.test)
			if err == nil {
				t.Errorf("expected error")
			}
		})
	}
}

// Get user
func TestAuthService_GetUser(t *testing.T) {
	service, userRepo, _ := setupAuthService()

	user, _ := models.NewUser("andrew", "andrewhoj@gmail.com", "12345678")
	userRepo.Save(user)

	lookupUser, err := service.GetUser(user.GetID().String())

	if err != nil {
		t.Errorf("expected error to be nil instead got: %s", err)
	}

	if lookupUser.GetID() != user.GetID() {
		t.Error("expected lookup users id to equal that of the mocked user")
	}
}

func TestAuthService_GetUserNotFound(t *testing.T) {
	service, _, _ := setupAuthService()
	_, err := service.GetUser("")

	if err == nil {
		t.Error("expected error, instead got nil")
	}
}

func TestAuthService_InvalidateRefreshToken(t *testing.T) {
	service, _, refreshTokenRepo := setupAuthService()

	user, _ := models.NewUser("andrew", "andrewhoj@gmail.com", "12345678")
	refreshToken, _ := models.NewRefreshToken(user.GetID())
	err := refreshTokenRepo.Save(refreshToken)

	if err != nil {
		t.Errorf("didn't expect error while saving refreshtoken but got: %s", err)
	}

	err = service.InvalidateRefreshToken(refreshToken.Token)
	if err != nil {
		t.Errorf("didn't expect error while invalidating refresh token: %s", err)
	}

	_, err = refreshTokenRepo.Find(refreshToken.Token)
	if err == nil {
		t.Error("expected error as the refresh token should no longer exist")
	}
}

func TestAuthService_InvalidateRefreshTokenInvalidToken(t *testing.T) {
	service, _, _ := setupAuthService()
	err := service.InvalidateRefreshToken("")

	if err == nil {
		t.Error("expected error as the given refresh token doesn't exist")
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	service, _, refreshTokenRepo := setupAuthService()

	user, _ := models.NewUser("andrew", "andrewhoj@gmail.com", "12345678")
	refreshToken, _ := models.NewRefreshToken(user.GetID())
	err := refreshTokenRepo.Save(refreshToken)
	if err != nil {
		t.Errorf("didn't expect any errors while saving the token, instead got: %s", err)
	}

	accessToken, err := service.RefreshAccessToken(refreshToken.Token)
	if err != nil {
		t.Errorf("didn't expect any errors while looking up refresh token, instead for: %s", err)
	}

	if len(accessToken) == 0 {
		t.Errorf("expected access token to have length above 0, instead it has: %v", len(accessToken))
	}
}

func TestAuthService_RefreshAccessTokenExpiredRefreshToken(t *testing.T) {
	service, _, refreshTokenRepo := setupAuthService()
	user, _ := models.NewUser("andrew", "andrewhoj@gmail.com", "12345678")
	token := &models.RefreshToken{
		ID:        "",
		UserID:    user.ID,
		Token:     "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(-time.Hour * 1),
	}
	refreshToken, err := models.CreateToken(token)
	token.Token = refreshToken

	err = refreshTokenRepo.Save(token)
	if err != nil {
		t.Errorf("didnt expect error while saving the refresh token, but got: %s", err)
	}

	if err != nil {
		t.Errorf("didn't expect error in creation of token, but got: %s", err)
	}

	_, err = service.RefreshAccessToken(refreshToken)
	if err == nil {
		t.Errorf("expected error while refreshing access token, instead got nil")
	}
}

func TestAuthService_RefreshAccessTokenNotFound(t *testing.T) {
	service, _, _ := setupAuthService()
	_, err := service.RefreshAccessToken("")

	if err == nil {
		t.Errorf("expected error while refreshing non-existing access token, instead got nil")
	}
}

func setupAuthService() (AuthService, repos.UserRepo, repos.RefreshTokenRepo) {
	logger := pgk.NewLogger()
	userRepo := &InMemoryUserRepo{users: make(map[string]*models.User)}
	tokenRepo := &InMemoryRefreshTokenRepo{tokens: make(map[string]*models.RefreshToken)}
	return NewAuthService(logger, userRepo, tokenRepo), userRepo, tokenRepo
}

type InMemoryUserRepo struct {
	users map[string]*models.User
}

func (u *InMemoryUserRepo) Save(user *models.User) error {
	u.users[user.ID.String()] = user
	return nil
}

func (u *InMemoryUserRepo) FindByEmail(email string) (*models.User, error) {
	for user := range u.users {
		if u.users[user].Email == email {
			return u.users[user], nil
		}
	}

	return nil, fmt.Errorf("no user with that email")
}

func (u *InMemoryUserRepo) FindByID(ID string) (*models.User, error) {
	user, ok := u.users[ID]
	if !ok {
		return nil, fmt.Errorf("no user with that id")
	}

	return user, nil
}

type InMemoryRefreshTokenRepo struct {
	tokens map[string]*models.RefreshToken
}

func (r *InMemoryRefreshTokenRepo) Save(token *models.RefreshToken) error {
	r.tokens[token.Token] = token
	return nil
}

func (r *InMemoryRefreshTokenRepo) Delete(tokenValue string) error {
	_, ok := r.tokens[tokenValue]
	if !ok {
		return fmt.Errorf("no token with that value")
	}

	delete(r.tokens, tokenValue)
	return nil
}

func (r *InMemoryRefreshTokenRepo) Find(tokenVal string) (*models.RefreshToken, error) {
	token, ok := r.tokens[tokenVal]
	if !ok {
		return nil, fmt.Errorf("no token with that value")
	}

	return token, nil
}

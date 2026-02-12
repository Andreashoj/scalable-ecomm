package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"fmt"
)

func SetupAuthService() (AuthService, repos.UserRepo, repos.RefreshTokenRepo) {
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

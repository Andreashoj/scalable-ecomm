package auth

import (
	"andreasho/scalable-ecomm/pgk"
	"andreasho/scalable-ecomm/services/user/internal/db/repos"
	"andreasho/scalable-ecomm/services/user/internal/domain"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func SetupAuthService(t *testing.T) (AuthService, repos.UserRepo, repos.RefreshTokenRepo) {
	//userRepo := &InMemoryUserRepo{users: make(map[string]*domain.User)}
	//tokenRepo := &InMemoryRefreshTokenRepo{tokens: make(map[string]*domain.RefreshToken)}

	db := pgk.SetupTestDB(t, "../services/user/internal/db/migrations")
	DB := sqlx.NewDb(db, "postgres")

	logger := pgk.NewLogger()
	userRepo := repos.NewUserRepo(DB)
	refreshTokenRepo := repos.NewRefreshTokenRepo(DB)

	return NewAuthService(logger, userRepo, refreshTokenRepo), userRepo, refreshTokenRepo
}

type InMemoryUserRepo struct {
	users map[string]*domain.User
}

func (u *InMemoryUserRepo) Save(user *domain.User) error {
	u.users[user.ID.String()] = user
	return nil
}

func (u *InMemoryUserRepo) FindByEmail(email string) (*domain.User, error) {
	for user := range u.users {
		if u.users[user].Email == email {
			return u.users[user], nil
		}
	}

	return nil, fmt.Errorf("no user with that email")
}

func (u *InMemoryUserRepo) FindByID(ID string) (*domain.User, error) {
	user, ok := u.users[ID]
	if !ok {
		return nil, fmt.Errorf("no user with that id")
	}

	return user, nil
}

type InMemoryRefreshTokenRepo struct {
	tokens map[string]*domain.RefreshToken
}

func (r *InMemoryRefreshTokenRepo) Save(token *domain.RefreshToken) error {
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

func (r *InMemoryRefreshTokenRepo) Find(tokenVal string) (*domain.RefreshToken, error) {
	token, ok := r.tokens[tokenVal]
	if !ok {
		return nil, fmt.Errorf("no token with that value")
	}

	return token, nil
}

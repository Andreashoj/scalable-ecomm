package auth

import "andreasho/scalable-ecomm/services/user/internal/db/models"

type mockRefreshTokenRepo struct {
}

func (r *mockRefreshTokenRepo) Save(token *models.RefreshToken) error {
	return nil
}

func (r *mockRefreshTokenRepo) Delete(refreshToken string) error {
	return nil
}

func (r *mockRefreshTokenRepo) Find(refreshToken string) (*models.RefreshToken, error) {
	return nil, nil
}

type mockUserRepo struct {
}

func (u *mockUserRepo) Save(user *models.User) error {
	return nil
}

func (u *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (u *mockUserRepo) FindByID(ID string) (*models.User, error) {
	return nil, nil
}

// Refresh access token

//

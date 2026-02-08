package repos

import (
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"database/sql"
	"fmt"
)

type RefreshTokenRepo interface {
	Save(token *models.RefreshToken) error
	Delete(refreshToken string) error
	Find(refreshToken string) (*models.RefreshToken, error)
}

type refreshTokenRepo struct {
	db *sql.DB
}

func (r *refreshTokenRepo) Find(refreshToken string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.QueryRow(`SELECT id, token, user_id, created_at, expires_at FROM refresh_token WHERE token = $1`, refreshToken).Scan(&token.ID, &token.Token, &token.UserID, &token.CreatedAt, &token.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving refresh token: %s", err)
	}

	return &token, nil
}

func (r *refreshTokenRepo) Delete(refreshToken string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_token WHERE token = $1`, refreshToken)
	if err != nil {
		return fmt.Errorf("failed deleting refresh token: %s", err)
	}

	return nil
}

func (r *refreshTokenRepo) Save(token *models.RefreshToken) error {
	_, err := r.db.Exec(`INSERT INTO refresh_token (id, token, expires_at, created_at, user_id) VALUES ($1, $2, $3, $4, $5)`, token.ID, token.Token, token.ExpiresAt, token.CreatedAt, token.UserID)
	if err != nil {
		return fmt.Errorf("failed creating refresh token in DB: %s", err)
	}

	return nil
}

func NewRefreshTokenRepo(db *sql.DB) RefreshTokenRepo {
	return &refreshTokenRepo{
		db: db,
	}
}

package repos

import (
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"database/sql"
	"fmt"
)

type RefreshTokenRepo interface {
	Save(token *models.RefreshToken) error
}

type refreshTokenRepo struct {
	db *sql.DB
}

func (r *refreshTokenRepo) Save(token *models.RefreshToken) error {
	_, err := r.db.Exec(`INSERT INTO refresh_token (id,  expires_at, created_at, user_id) VALUES ($1, $2, $3, $4)`, token.ID, token.ExpiresAt, token.CreatedAt, token.UserID)
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

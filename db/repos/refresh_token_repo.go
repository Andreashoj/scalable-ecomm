package repos

import (
	"andreasho/scalable-ecomm/db/models"
)

type RefreshTokenRepo interface {
	Save(token *models.RefreshToken) error
}

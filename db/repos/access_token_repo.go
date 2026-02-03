package repos

import "andreasho/scalable-ecomm/db/models"

type AccessTokenRepo interface {
	Save(token *models.AccessToken) error
}

type accessTokenRepo struct{}

func (a accessTokenRepo) Save(token *models.AccessToken) error {
	//TODO implement me
	panic("implement me")
}

func NewAccessTokenRepo() AccessTokenRepo {
	return &accessTokenRepo{}
}

package repos

import "andreasho/scalable-ecomm/db/models"

type UserRepo interface {
	Save(user *models.User) error
	Find(email string) (*models.User, error)
}

type userRepo struct{}

func (u *userRepo) Find(email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) Save(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

package services

type UserService interface {
	IsAdmin(accessToken string) bool
}

type userService struct {
}

func (u *userService) IsAdmin(accessToken string) bool {
	return true
}

func NewUserService() UserService {
	return &userService{}
}

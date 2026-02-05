package repos

import (
	"andreasho/scalable-ecomm/services/user/internal/db/models"
	"database/sql"
	"fmt"
)

type UserRepo interface {
	Save(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(ID string) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func (u *userRepo) FindByID(ID string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(`SELECT id, name, email, role, created_at FROM users WHERE id = $1`, ID).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by email: %s", err)
	}

	return &user, nil
}

func (u *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(`SELECT id, name, email, role, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by email: %s", err)
	}

	return &user, nil
}

func (u *userRepo) Save(user *models.User) error {
	_, err := u.db.Exec(
		`INSERT INTO users (id, name, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed creating user in DB: %s", err)
	}

	return nil
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

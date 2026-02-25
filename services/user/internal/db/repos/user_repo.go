package repos

import (
	"andreasho/scalable-ecomm/services/user/internal/domain"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Save(user *domain.User) error
	Update(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(ID string) (*domain.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func (u *userRepo) Update(user *domain.User) error {
	_, err := squirrel.
		Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(u.db).
		Exec()

	return err
}

func (u *userRepo) FindByID(ID string) (*domain.User, error) {
	var user domain.User
	err := u.db.QueryRow(`SELECT id, name, email, role, created_at FROM users WHERE id = $1`, ID).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by email: %s", err)
	}

	return &user, nil
}

func (u *userRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := u.db.QueryRow(`SELECT id, name, email, role, password, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by email: %s", err)
	}

	return &user, nil
}

func (u *userRepo) Save(user *domain.User) error {
	_, err := u.db.Exec(
		`INSERT INTO users (id, name, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("email already exists")
		}

		return fmt.Errorf("failed creating user in DB: %s", err)
	}

	return nil
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

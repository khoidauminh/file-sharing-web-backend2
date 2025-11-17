package repository

import (
	"database/sql"
	"fmt"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(DB *sql.DB) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (ur *SQLUserRepository) Create(user *domain.User) error {
	row := ur.db.QueryRow("INSERT INTO users (name, password, Email, Role, enableTOTP) VALUES ($1, $2, $3, $4, $5) RETURNING user_id", user.Username, user.Password, user.Email, user.Role, user.EnableTOTP)
	err := row.Scan(&user.Id)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *SQLUserRepository) FindById(id int, user *domain.User) error {
	row := ur.db.QueryRow("SELECT * FROM users WHERE user_id = $1", id)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Role, &user.EnableTOTP)

	if err != nil {
		return err
	}

	return nil
}

func (ur *SQLUserRepository) FindByEmail(email string, user *domain.User) error {
	row := ur.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Role, &user.EnableTOTP)
	if err != nil {
		return err
	}

	return nil
}

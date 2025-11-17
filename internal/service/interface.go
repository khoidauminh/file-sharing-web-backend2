package service

import "github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"

type UserService interface {
	CreateUser(username, password, email, role string) error
	GetUserById(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

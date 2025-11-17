package repository

import "github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindById(id int, user *domain.User) error
	FindByEmail(email string, user *domain.User) error
}

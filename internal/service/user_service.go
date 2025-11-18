package service

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}

func (us *userService) GetUserById(id string) (*domain.User, error) {
	user := &domain.User{}
	err := us.userRepo.FindById(id, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := us.userRepo.FindByEmail(email, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

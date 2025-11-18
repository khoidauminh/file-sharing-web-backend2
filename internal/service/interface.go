package service

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUserById(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type AuthService interface {
	CreateUser(username, password, email, role string) (*domain.User, error)
	Login(email, password string) (user *domain.User, accessToken string, expiresIn int, err error)
	Logout(ctx *gin.Context) error
}

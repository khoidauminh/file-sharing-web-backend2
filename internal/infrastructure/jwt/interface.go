package jwt

import "github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"

type TokenService interface {
	GenerateAccessToken(user domain.User) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

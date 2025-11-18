package app

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/handlers"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/routes"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/infrastructure/jwt"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/repository"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService jwt.TokenService) *AuthModule {
	userRepository := repository.NewSQLUserRepository(ctx.DB)
	authRepository := repository.NewAuthRepository(ctx.DB)
	authService := service.NewAuthService(userRepository, authRepository, tokenService)
	authHandler := handlers.NewAuthHandler(authService)
	authRoutes := routes.NewAuthRoutes(authHandler)
	return &AuthModule{routes: authRoutes}
}

func (m *AuthModule) Routes() routes.Route {
	return m.routes
}

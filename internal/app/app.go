package app

import (
	"database/sql"
	"log"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/config"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/routes"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB *sql.DB
}

func NewApplication(cfg *config.Config) *Application {

	r := gin.Default()

	if err := database.InitDB(); err != nil {
		log.Fatal("unable to connnect to db")
	}

	ctx := &ModuleContext{
		DB: database.DB,
	}

	modules := []Module{
		NewUserModule(ctx),
	}

	routes.RegisterRoutes(r, getModulRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	if a.config.ServerAddress == "" {
		a.config.ServerAddress = ":8080"
	}

	log.Printf(" Server is running at http://localhost%s\n", a.config.ServerAddress)
	return a.router.Run(a.config.ServerAddress)
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

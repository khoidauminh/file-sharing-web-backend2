// file: internal/app/admin_module.go (Mới)
package app

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/config"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/handlers"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/routes"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
)

type adminModule struct {
	routes routes.Route
}

// Giả định NewAdminService nhận *config.Config (chứa Policy tĩnh)
func NewAdminModule(cfg *config.Config) Module {

	// Policy tĩnh: không cần Repository
	adminService := service.NewAdminService(cfg)
	adminHandler := handlers.NewAdminHandler(adminService)
	adminRoutes := routes.NewAdminRoutes(adminHandler)

	return &adminModule{
		routes: adminRoutes,
	}
}

func (m *adminModule) Routes() routes.Route {
	return m.routes
}

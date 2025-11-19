package routes

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	handler *handlers.AdminHandler
}

func NewAdminRoutes(handler *handlers.AdminHandler) *AdminRoutes {
	return &AdminRoutes{
		handler: handler,
	}
}

func (ar *AdminRoutes) Register(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		// Cần có middleware kiểm tra quyền Admin tại đây
		admin.GET("/policy", ar.handler.GetSystemPolicy)      // Lấy cấu hình hệ thống
		admin.PATCH("/policy", ar.handler.UpdateSystemPolicy) // Cập nhật cấu hình hệ thống

		// Cleanup có thể là protected route cho admin/cron job
		admin.POST("/cleanup", ar.handler.CleanupExpiredFiles) // Xóa file hết hạn
	}
}

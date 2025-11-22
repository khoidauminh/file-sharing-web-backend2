package routes

import (
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/handlers"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/middleware"
	"github.com/gin-gonic/gin"
)

type FileRoutes struct {
	handler *handlers.FileHandler
}

func NewFileRoutes(handler *handlers.FileHandler) *FileRoutes {
	return &FileRoutes{
		handler: handler,
	}
}

func (fr *FileRoutes) Register(r *gin.RouterGroup) {
	files := r.Group("/files")
	optional := files.Group("/")
	optional.Use(middleware.AuthMiddlewareUpload())
	{
		// POST /api/files/upload
		// Lưu ý: Endpoint này không cần AuthMiddleware() nếu là Anonymous Upload.
		// Cần middleware kiểm tra JWT VÀ cho phép request tiếp tục nếu không có token.
		// Hiện tại nó đang được đăng ký dưới protected group, nhưng logic xử lý trong handler đã cho phép anonymous.
		optional.POST("/upload", fr.handler.UploadFile)
	}
	protected := files.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// GET /api/files/my
		protected.GET("/my", fr.handler.GetMyFiles)

		// DELETE /api/files/:id
		protected.DELETE("/:id", fr.handler.DeleteFile)

		// Các routes download công khai và xem thông tin (chưa triển khai đầy đủ)
		protected.GET("/:shareToken", fr.handler.GetFileInfo)
		protected.GET("/:shareToken/download", fr.handler.DownloadFile)
	}
}

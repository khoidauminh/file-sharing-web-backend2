package handlers

import (
	"net/http"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/dto"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/pkg/utils"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/pkg/validation"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	admin_service service.AdminService
}

func NewAdminHandler(admin_service service.AdminService) *AdminHandler {
	return &AdminHandler{
		admin_service: admin_service,
	}
}

func (ah *AdminHandler) GetSystemPolicy(ctx *gin.Context) {
	policy, err := ah.admin_service.GetSystemPolicy(ctx)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	// Trả về cấu hình hệ thống
	utils.ResponseSuccess(ctx, http.StatusOK, "System policy retrieved successfully", policy)
}

func (ah *AdminHandler) UpdateSystemPolicy(ctx *gin.Context) {
	var req dto.UpdatePolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	updatedPolicy, err := ah.admin_service.UpdateSystemPolicy(ctx, req.ToMap())

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	// Cập nhật thành công
	utils.ResponseSuccess(ctx, http.StatusOK, "System policy updated successfully", gin.H{"policy": updatedPolicy})
}

func (ah *AdminHandler) CleanupExpiredFiles(ctx *gin.Context) {
	// Endpoint này có thể được gọi bằng Admin token hoặc X-Cron-Secret
	// Giả định Middleware đã xác thực quyền Admin hoặc Secret

	deletedCount, err := ah.admin_service.CleanupExpiredFiles(ctx)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Cleanup completed", gin.H{
		"deletedFiles": deletedCount,
	})
}

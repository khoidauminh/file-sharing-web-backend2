package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/api/dto"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/pkg/utils"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	file_service service.FileService
}

func NewFileHandler(file_service service.FileService) *FileHandler {
	return &FileHandler{
		file_service: file_service,
	}
}

func (fh *FileHandler) UploadFile(ctx *gin.Context) {
	var req dto.UploadRequest

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("File is required for upload", utils.ErrCodeBadRequest))
		return
	}

	req.FileNameForValidation = fileHeader.Filename

	if err := ctx.ShouldBind(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var userID *string
	if val, exists := ctx.Get("userID"); exists && val != "" {
		strVal := val.(string)
		userID = &strVal
	} else {
		userID = nil
	}

	uploadedFile, err := fh.file_service.UploadFile(ctx, fileHeader, &req, userID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := gin.H{
		"id":            uploadedFile.Id,
		"fileName":      uploadedFile.FileName,
		"fileSize":      uploadedFile.FileSize,
		"shareToken":    uploadedFile.ShareToken,
		"shareLink":     fmt.Sprintf("http://%s/f/%s", ctx.Request.Host, uploadedFile.ShareToken),
		"isPublic":      uploadedFile.IsPublic,
		"hasPassword":   uploadedFile.HasPassword,
		"availableFrom": uploadedFile.AvailableFrom,
		"availableTo":   uploadedFile.AvailableTo,
		"validityDays":  uploadedFile.ValidityDays,
		"enableTOTP":    uploadedFile.EnableTOTP,
		"createdAt":     uploadedFile.CreatedAt,
	}

	//utils.ResponseSuccess(ctx, http.StatusCreated, "File uploaded successfully", gin.H{"file": response})
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"file":    response,
		"message": "File uploaded successfully",
	})
}

func (fh *FileHandler) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Param("id")

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("Unauthorized access", utils.ErrCodeUnauthorized))
		return
	}

	err := fh.file_service.DeleteFile(ctx, fileID, userID.(string))

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	//utils.ResponseSuccess(ctx, http.StatusOK, "File deleted successfully", gin.H{"fileId": fileID})
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
		"fileId":  fileID,
	})
}

func (fh *FileHandler) GetMyFiles(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("Unauthorized access", utils.ErrCodeUnauthorized))
		return
	}

	status := ctx.DefaultQuery("status", "all")
	page := utils.GetIntQuery(ctx, "page", 1)
	limit := utils.GetIntQuery(ctx, "limit", 20)
	sortBy := ctx.DefaultQuery("sortBy", "createdAt")
	order := ctx.DefaultQuery("order", "desc")

	params := domain.ListFileParams{
		Status: strings.ToLower(status),
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  strings.ToLower(order),
	}

	result, err := fh.file_service.GetMyFiles(ctx, userID.(string), params)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	//utils.ResponseSuccess(ctx, http.StatusOK, "User files retrieved successfully", gin.H{"file": result})
	ctx.JSON(http.StatusOK, gin.H{
		"file" : result,
	})
}

func (fh *FileHandler) GetFileInfo(ctx *gin.Context) {
	ident := ctx.Param("ident")
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = nil
	}

	var result *domain.File = nil
	var err error = nil

	if uuid.Validate(ident) == nil {
		result, err = fh.file_service.GetFileInfoID(ctx, ident, userID.(string))
	} else {
		result, err = fh.file_service.GetFileInfo(ctx, ident, userID.(string))
	}

	if err != nil {
		utils.ResponseError(ctx, utils.WrapError(err, "Failed to access file", utils.ErrCodeInternal))
		return
	}

	//utils.ResponseSuccess(ctx, http.StatusOK, "File retrieved successfully", gin.H{"file": result})
	ctx.JSON(http.StatusOK, gin.H{
		"file" : result,
	})
}

func (fh *FileHandler) DownloadFile(ctx *gin.Context) {
	fileToken := ctx.Param("ident")
	password := ctx.Query("password")
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = nil
	}

	info, file, err := fh.file_service.DownloadFile(ctx, fileToken, userID.(string), password)
	if err != nil {
		utils.ResponseError(ctx, utils.WrapError(err, "Failed to download file", utils.ErrCodeInternal))
		return
	}

	ctx.Data(http.StatusOK, info.MimeType, file)
}

func (fh *FileHandler) GetFileDownloadHistory(ctx *gin.Context) {
	fileID := ctx.Param("ident")
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = nil
	}

	page := utils.GetIntQuery(ctx, "page", 1)
	limit := utils.GetIntQuery(ctx, "limit", 20)
	if limit == 0 {
		utils.ResponseError(ctx, utils.NewError("Limit must not be 0", utils.ErrCodeBadRequest))
		return
	}

	history, err := fh.file_service.GetFileDownloadHistory(ctx, fileID, userID.(string), page, limit)
	if err != nil {
		utils.ResponseError(ctx, utils.WrapError(err, "Failed to get file history", utils.ErrCodeInternal))
		return
	}

	ctx.JSON(http.StatusOK, history)
}

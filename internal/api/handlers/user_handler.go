package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user_service service.UserService
}

func NewUserHandler(user_service service.UserService) *UserHandler {
	return &UserHandler{
		user_service: user_service,
	}
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.user_service.CreateUser(user.Username, user.Password, user.Email, user.Role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

func (uh *UserHandler) GetUserById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user domain.User
	createdUser, err := uh.user_service.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	} else {
		user = *createdUser
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (uh *UserHandler) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}
	var user domain.User
	createdUser, err := uh.user_service.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	} else {
		user = *createdUser
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

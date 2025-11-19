package handlers

import (
	"net/http"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/infrastructure/jwt"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/service"
	"github.com/gin-gonic/gin"
)

type TotpHandler struct {
	TotpService service.TotpService
}

func NewTotpHandler(totpService service.TotpService) *TotpHandler {
	return &TotpHandler{
		TotpService: totpService,
	}
}

func getUserIDFromContext(c *gin.Context) (string, bool) {
	userObj, exists := c.Get("user")
	if !exists {
		return "", false
	}
	
	claims, ok := userObj.(*jwt.Claims)
	if !ok {
		return "", false
	}
	
	return claims.UserID, true
}

func (h *TotpHandler) SetupTOTP(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized or invalid token"})
		return
	}

	resp, err := h.TotpService.SetupTOTP(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "TOTP secret generated",
		"totpSetup": resp,
	})
}

type VerifyTOTPRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *TotpHandler) VerifyTOTP(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized or invalid token"})
		return
	}

	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	okVerify, err := h.TotpService.VerifyTOTP(userID, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !okVerify {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "TOTP verified successfully",
		"totpEnabled": true,
	})
}
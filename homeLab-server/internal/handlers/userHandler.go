package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var credentials entities.UserCredentials // Ensure this struct is defined correctly
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.userUseCase.Login(credentials.Username, credentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	// Assuming userID is passed in the context or headers
	userID := ctx.Param("userID") // Adjust as necessary
	if err := h.userUseCase.Logout(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *UserHandler) ViewUserRoles(ctx *gin.Context) {
	roles, err := h.userUseCase.ViewUserRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}

func (h *UserHandler) ViewAuditLogs(ctx *gin.Context) {
	logs, err := h.userUseCase.ViewAuditLogs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}

func (h *UserHandler) GetPermissions(c *gin.Context) {
	// Handle fetching user permissions
}

func (h *UserHandler) GetAuditLogs(c *gin.Context) {
	// Handle fetching audit logs
}

func (h *AuthenticationHandler) Register(c *gin.Context) {
	var credentials entities.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authenticationUseCase.Register(credentials)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"x-auth-token": token})
}

func (h *AuthenticationHandler) Login(c *gin.Context) {
	var credentials entities.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authenticationUseCase.Login(credentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"x-auth-token": token})
}

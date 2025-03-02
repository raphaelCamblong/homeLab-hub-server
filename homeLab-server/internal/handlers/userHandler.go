package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var credentials Credentials // Define Credentials struct as per your requirements
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.userUseCase.Login(credentials)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	if err := h.userUseCase.Logout(); err != nil {
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
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type HealthHandler struct {
	statusUseCase usecase.HealthStatusUseCase
}

func NewHealthHandler(statusUseCase usecase.HealthStatusUseCase) *HealthHandler {
	return &HealthHandler{statusUseCase: statusUseCase}
}

func (h *HealthHandler) GetStatus(ctx *gin.Context) {
	status, err := h.statusUseCase.GetStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func (h *HealthHandler) GetOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server Ok"})
}

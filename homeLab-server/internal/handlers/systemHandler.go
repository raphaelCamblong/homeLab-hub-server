package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type SystemHandler struct {
	systemUseCase usecase.SystemUseCase
}

func NewSystemHandler(systemUseCase usecase.SystemUseCase) *SystemHandler {
	return &SystemHandler{systemUseCase: systemUseCase}
}

func (h *SystemHandler) GetOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server Ok"})
}

func (h *SystemHandler) GetSystemInfo(ctx *gin.Context) {
	info, err := h.systemUseCase.GetSystemInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func (h *SystemHandler) CheckHealth(ctx *gin.Context) {
	health, err := h.systemUseCase.CheckHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, health)
}

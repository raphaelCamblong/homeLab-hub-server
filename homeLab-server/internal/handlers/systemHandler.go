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

func (h *SystemHandler) GetSystemInfo(ctx *gin.Context) {
	info, err := h.systemUseCase.GetSystemInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func (h *SystemHandler) ListProcesses(ctx *gin.Context) {
	processes, err := h.systemUseCase.ListProcesses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, processes)
}

func (h *SystemHandler) ListActiveServices(ctx *gin.Context) {
	services, err := h.systemUseCase.ListActiveServices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, services)
}

func (h *SystemHandler) FetchSystemMetrics(ctx *gin.Context) {
	metrics, err := h.systemUseCase.FetchSystemMetrics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}

func (h *SystemHandler) CheckHealth(ctx *gin.Context) {
	health, err := h.systemUseCase.CheckHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, health)
}

func (h *SystemHandler) RestartSystem(ctx *gin.Context) {
	if err := h.systemUseCase.RestartSystem(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "System restarted successfully"})
}

func (h *SystemHandler) ShutdownSystem(ctx *gin.Context) {
	if err := h.systemUseCase.ShutdownSystem(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "System shutdown successfully"})
} 
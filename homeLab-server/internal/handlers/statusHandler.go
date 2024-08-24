package handlers

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
	"net/http"
)

type StatusHandler struct {
	statusUseCase usecase.StatusUseCase
}

func NewStatusHandler(statusUseCase usecase.StatusUseCase) *StatusHandler {
	return &StatusHandler{statusUseCase: statusUseCase}
}

func (h *StatusHandler) GetStatus(ctx *gin.Context) {
	status, err := h.statusUseCase.GetStatus(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

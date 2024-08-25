package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type ThermalHandler struct {
	thermalUseCase usecase.ThermalUseCase
}

func NewThermalHandler(thermalUseCase usecase.ThermalUseCase) *ThermalHandler {
	return &ThermalHandler{thermalUseCase: thermalUseCase}
}

func (h *ThermalHandler) GetThermal(ctx *gin.Context) {
	thermal, err := h.thermalUseCase.GetThermal(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, thermal)
}

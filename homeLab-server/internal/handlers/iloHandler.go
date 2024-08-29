package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type IloHandler struct {
	thermalUseCase usecase.IlohUseCase
}

func NewIloHandler(thermalUseCase usecase.IlohUseCase) *IloHandler {
	return &IloHandler{thermalUseCase: thermalUseCase}
}

func (h *IloHandler) GetThermal(ctx *gin.Context) {
	thermal, err := h.thermalUseCase.GetThermal()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, thermal)
}

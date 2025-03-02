package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type MonitoringHandler struct {
	monitoringUseCase usecase.MonitoringUseCase
}

func NewMonitoringHandler(monitoringUseCase usecase.MonitoringUseCase) *MonitoringHandler {
	return &MonitoringHandler{monitoringUseCase: monitoringUseCase}
}

func (h *MonitoringHandler) QueryPrometheus(ctx *gin.Context) {
	data, err := h.monitoringUseCase.QueryPrometheus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *MonitoringHandler) GetGrafanaDashboardData(ctx *gin.Context) {
	data, err := h.monitoringUseCase.GetGrafanaDashboardData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *MonitoringHandler) FetchActiveAlerts(ctx *gin.Context) {
	alerts, err := h.monitoringUseCase.FetchActiveAlerts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, alerts)
} 
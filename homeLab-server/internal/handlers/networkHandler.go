package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type NetworkHandler struct {
	networkUseCase usecase.NetworkUseCase
}

func NewNetworkHandler(networkUseCase usecase.NetworkUseCase) *NetworkHandler {
	return &NetworkHandler{networkUseCase: networkUseCase}
}

func (h *NetworkHandler) ListDevices(ctx *gin.Context) {
	devices, err := h.networkUseCase.ListDevices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, devices)
}

func (h *NetworkHandler) ListNetworkInterfaces(ctx *gin.Context) {
	interfaces, err := h.networkUseCase.ListNetworkInterfaces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, interfaces)
}

func (h *NetworkHandler) GetInterfaceDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	details, err := h.networkUseCase.GetInterfaceDetails(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, details)
}

func (h *NetworkHandler) ViewFirewallRules(ctx *gin.Context) {
	rules, err := h.networkUseCase.ViewFirewallRules()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rules)
}

func (h *NetworkHandler) ModifyFirewallRules(ctx *gin.Context) {
	var rules []Rule // Define Rule struct as per your requirements
	if err := ctx.ShouldBindJSON(&rules); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.networkUseCase.ModifyFirewallRules(rules); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Firewall rules modified successfully"})
}

func (h *NetworkHandler) ViewDNSSettings(ctx *gin.Context) {
	settings, err := h.networkUseCase.ViewDNSSettings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (h *NetworkHandler) RunNetworkSpeedTest(ctx *gin.Context) {
	result, err := h.networkUseCase.RunNetworkSpeedTest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
} 
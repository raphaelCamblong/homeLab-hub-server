package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type K3sHandler struct {
	k3sUseCase usecase.K3sUseCase
}

func NewK3sHandler(k3sUseCase usecase.K3sUseCase) *K3sHandler {
	return &K3sHandler{k3sUseCase: k3sUseCase}
}

func (h *K3sHandler) ListK3sNodes(ctx *gin.Context) {
	nodes, err := h.k3sUseCase.ListK3sNodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (h *K3sHandler) ListRunningPods(ctx *gin.Context) {
	pods, err := h.k3sUseCase.ListRunningPods()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pods)
}

func (h *K3sHandler) ListPodsInNamespace(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	pods, err := h.k3sUseCase.ListPodsInNamespace(namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pods)
}

func (h *K3sHandler) ListK3sServices(ctx *gin.Context) {
	services, err := h.k3sUseCase.ListK3sServices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, services)
}

func (h *K3sHandler) ListK3sDeployments(ctx *gin.Context) {
	deployments, err := h.k3sUseCase.ListK3sDeployments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deployments)
}

func (h *K3sHandler) GetPodLogs(ctx *gin.Context) {
	pod := ctx.Param("pod")
	logs, err := h.k3sUseCase.GetPodLogs(pod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}

func (h *K3sHandler) GetClusterEvents(ctx *gin.Context) {
	events, err := h.k3sUseCase.GetClusterEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func (h *K3sHandler) ApplyYAMLManifest(ctx *gin.Context) {
	var manifest string // Define manifest as per your requirements
	if err := ctx.ShouldBindJSON(&manifest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.k3sUseCase.ApplyYAMLManifest(manifest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "YAML manifest applied successfully"})
}

func (h *K3sHandler) DeleteResource(ctx *gin.Context) {
	resourceID := ctx.Query("id")
	if err := h.k3sUseCase.DeleteResource(resourceID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
} 
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type ClusterHandler struct {
	ClustersUseCase usecase.ClusterUseCase
}

func NewClusterHandler(ClustersUseCase usecase.ClusterUseCase) *ClusterHandler {
	return &ClusterHandler{ClustersUseCase: ClustersUseCase}
}

func (h *ClusterHandler) ListClusterNodes(ctx *gin.Context) {
	nodes, err := h.ClustersUseCase.ListClusterNodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (h *ClusterHandler) ListRunningPods(ctx *gin.Context) {
	pods, err := h.ClustersUseCase.ListRunningPods()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pods)
}

func (h *ClusterHandler) ListPodsInNamespace(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	pods, err := h.ClustersUseCase.ListPodsInNamespace(namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pods)
}

func (h *ClusterHandler) ListClusterServices(ctx *gin.Context) {
	services, err := h.ClustersUseCase.ListClusterServices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, services)
}

func (h *ClusterHandler) ListClusterDeployments(ctx *gin.Context) {
	deployments, err := h.ClustersUseCase.ListClusterDeployments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deployments)
}

func (h *ClusterHandler) GetPodLogs(ctx *gin.Context) {
	pod := ctx.Param("pod")
	logs, err := h.ClustersUseCase.GetPodLogs(pod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}

func (h *ClusterHandler) GetClusterEvents(ctx *gin.Context) {
	events, err := h.ClustersUseCase.GetClusterEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func (h *ClusterHandler) ApplyYAMLManifest(ctx *gin.Context) {
	var manifest string // Define manifest as per your requirements
	if err := ctx.ShouldBindJSON(&manifest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.ClustersUseCase.ApplyYAMLManifest(manifest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "YAML manifest applied successfully"})
}

func (h *ClusterHandler) DeleteResource(ctx *gin.Context) {
	resourceID := ctx.Query("id")
	if err := h.ClustersUseCase.DeleteResource(resourceID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
} 
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type ClusterHandler struct {
	clusterUseCase usecase.ClusterUseCase
}

func NewClusterHandler(clusterUseCase usecase.ClusterUseCase) *ClusterHandler {
	return &ClusterHandler{
		clusterUseCase: clusterUseCase,
	}
}

// GetClusterNodes godoc
// @Summary List all cluster nodes
// @Description Get information about all nodes in the Kubernetes cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.ClusterNode
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/nodes [get]
func (h *ClusterHandler) GetClusterNodes(c *gin.Context) {
	nodes, err := h.clusterUseCase.GetClusterNodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// GetNodeStats godoc
// @Summary Get node statistics
// @Description Get detailed statistics for a specific node in the cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Param nodeID path string true "Node ID"
// @Security BearerAuth
// @Success 200 {object} entities.NodeStats
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/nodes/{nodeID}/stats [get]
func (h *ClusterHandler) GetNodeStats(c *gin.Context) {
	nodeID := c.Param("nodeID")
	stats, err := h.clusterUseCase.GetNodeStats(nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetRunningPods godoc
// @Summary List all running pods
// @Description Get information about all running pods across all namespaces
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Pod
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/pods [get]
func (h *ClusterHandler) GetRunningPods(c *gin.Context) {
	pods, err := h.clusterUseCase.GetRunningPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pods)
}

// GetWorkloads godoc
// @Summary List all workflows
// @Description Get information about all workflows in the cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Workflow
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/workflows [get]
func (h *ClusterHandler) GetWorkloads(c *gin.Context) {
	var workflows []entities.Workflow
	workflows, err := h.clusterUseCase.GetWorkloads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workflows)
}

// GetPodsInNamespace godoc
// @Summary List pods in namespace
// @Description Get information about all pods in a specific namespace
// @Tags cluster
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace name"
// @Security BearerAuth
// @Success 200 {array} entities.Pod
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/pods/{namespace} [get]
func (h *ClusterHandler) GetPodsInNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	pods, err := h.clusterUseCase.GetPodsInNamespace(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pods)
}

// GetClusterServices godoc
// @Summary List all services
// @Description Get information about all services in the cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Service
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/services [get]
func (h *ClusterHandler) GetClusterServices(c *gin.Context) {
	services, err := h.clusterUseCase.GetClusterServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetClusterDeployments godoc
// @Summary List all deployments
// @Description Get information about all deployments in the cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Deployment
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/deployments [get]
func (h *ClusterHandler) GetClusterDeployments(c *gin.Context) {
	deployments, err := h.clusterUseCase.GetClusterDeployments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deployments)
}

// GetClusterStats godoc
// @Summary Get cluster statistics
// @Description Get statistics about nodes, pods, deployments, and statefulsets in the cluster
// @Tags cluster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ClusterStats
// @Failure 500 {object} gin.H
// @Router /api/v1/cluster/ [get]
func (h *ClusterHandler) GetClusterStats(c *gin.Context) {
	stats, err := h.clusterUseCase.GetClusterStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

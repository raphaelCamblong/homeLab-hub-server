package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ClusterRepository interface {
	ListClusterNodes() ([]entities.ClusterNode, error)
	GetNodeStats(nodeID string) (entities.NodeStats, error)
	ListRunningPods() ([]entities.Pod, error)
	ListWorkflows() ([]entities.Workflow, error)
	ListPodsInNamespace(namespace string) ([]entities.Pod, error)
	ListClusterServices() ([]entities.Service, error)
	ListClusterDeployments() ([]entities.Deployment, error)
	GetClusterStats() (*entities.ClusterStats, error)
}

type ClustersRepository struct {
	k8sClient *client.K8sClient
}

func NewClusterRepository(k8sClient *client.K8sClient) ClusterRepository {
	return &ClustersRepository{
		k8sClient: k8sClient,
	}
}

func (r *ClustersRepository) ListClusterNodes() ([]entities.ClusterNode, error) {
	return r.k8sClient.ListNodes()
}

func (r *ClustersRepository) GetNodeStats(nodeID string) (entities.NodeStats, error) {
	return r.k8sClient.GetNodeStats(nodeID)
}

func (r *ClustersRepository) ListRunningPods() ([]entities.Pod, error) {
	return r.k8sClient.ListPods("")
}

func (r *ClustersRepository) ListWorkflows() ([]entities.Workflow, error) {
	// Note: This is a placeholder. In a real implementation, you would need to
	// implement workflow listing based on your workflow engine (e.g., Argo, Tekton)
	return []entities.Workflow{}, nil
}

func (r *ClustersRepository) ListPodsInNamespace(namespace string) ([]entities.Pod, error) {
	return r.k8sClient.ListPods(namespace)
}

func (r *ClustersRepository) ListClusterServices() ([]entities.Service, error) {
	return r.k8sClient.ListServices("")
}

func (r *ClustersRepository) ListClusterDeployments() ([]entities.Deployment, error) {
	return r.k8sClient.ListDeployments("")
}

func (r *ClustersRepository) GetClusterStats() (*entities.ClusterStats, error) {
	return r.k8sClient.GetClusterStats()
}

package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ClusterUseCase interface {
	GetClusterNodes() ([]entities.ClusterNode, error)
	GetNodeStats(nodeID string) (entities.NodeStats, error)
	GetRunningPods() ([]entities.Pod, error)
	GetWorkloads() ([]entities.Workflow, error)
	GetPodsInNamespace(namespace string) ([]entities.Pod, error)
	GetClusterServices() ([]entities.Service, error)
	GetClusterDeployments() ([]entities.Deployment, error)
	GetClusterStats() (*entities.ClusterStats, error)
}

type clusterUseCase struct {
	clusterRepo repositories.ClusterRepository
}

func NewClusterUseCase(clusterRepo repositories.ClusterRepository) ClusterUseCase {
	return &clusterUseCase{
		clusterRepo: clusterRepo,
	}
}

func (u *clusterUseCase) GetClusterNodes() ([]entities.ClusterNode, error) {
	return u.clusterRepo.ListClusterNodes()
}

func (u *clusterUseCase) GetNodeStats(nodeID string) (entities.NodeStats, error) {
	return u.clusterRepo.GetNodeStats(nodeID)
}

func (u *clusterUseCase) GetRunningPods() ([]entities.Pod, error) {
	return u.clusterRepo.ListRunningPods()
}

func (u *clusterUseCase) GetWorkloads() ([]entities.Workflow, error) {
	return u.clusterRepo.ListWorkflows()
}

func (u *clusterUseCase) GetPodsInNamespace(namespace string) ([]entities.Pod, error) {
	return u.clusterRepo.ListPodsInNamespace(namespace)
}

func (u *clusterUseCase) GetClusterServices() ([]entities.Service, error) {
	return u.clusterRepo.ListClusterServices()
}

func (u *clusterUseCase) GetClusterDeployments() ([]entities.Deployment, error) {
	return u.clusterRepo.ListClusterDeployments()
}

func (u *clusterUseCase) GetClusterStats() (*entities.ClusterStats, error) {
	return u.clusterRepo.GetClusterStats()
}

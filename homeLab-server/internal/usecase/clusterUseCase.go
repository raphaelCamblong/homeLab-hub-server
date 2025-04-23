package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ClusterUseCase interface {
	ListNodes() ([]ClusterNode, error)
	ListRunningPods() ([]Pod, error)
	ListPodsInNamespace(namespace string) ([]Pod, error)
	ListServices() ([]Service, error)
	ListDeployments() ([]Deployment, error)
	GetPodLogs(pod string) (string, error)
	GetEvents() ([]Event, error)
}

type ClustersUseCase struct {
	ClustersRepository repositories.ClusterRepository
}

func NewClusterUseCase(ClustersRepository repositories.ClusterRepository) ClusterUseCase {
	return &ClustersUseCase{ClustersRepository: ClustersRepository}
}

func (u *ClustersUseCase) ListNodes() ([]ClusterNode, error) {
	return u.ClustersRepository.ListClusterNodes()
}

func (u *ClustersUseCase) ListRunningPods() ([]Pod, error) {
	return u.ClustersRepository.ListRunningPods()
}

func (u *ClustersUseCase) ListPodsInNamespace(namespace string) ([]Pod, error) {
	return u.ClustersRepository.ListPodsInNamespace(namespace)
}

func (u *ClustersUseCase) ListServices() ([]Service, error) {
	return u.ClustersRepository.ListServices()
}

func (u *ClustersUseCase) ListDeployments() ([]Deployment, error) {
	return u.ClustersRepository.ListDeployments()
}

func (u *ClustersUseCase) GetPodLogs(pod string) (string, error) {
	return u.ClustersRepository.GetPodLogs(pod)
}

func (u *ClustersUseCase) GetEvents() ([]Event, error) {
	return u.ClustersRepository.GetClusterEvents()
}

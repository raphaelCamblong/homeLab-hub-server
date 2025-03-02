package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type K3sUseCase interface {
	ListK3sNodes() ([]K3sNode, error)
	ListRunningPods() ([]Pod, error)
	ListPodsInNamespace(namespace string) ([]Pod, error)
	ListK3sServices() ([]Service, error)
	ListK3sDeployments() ([]Deployment, error)
	GetPodLogs(pod string) (string, error)
	GetClusterEvents() ([]Event, error)
	ApplyYAMLManifest(manifest string) error
	DeleteResource(resourceID string) error
}

type k3sUseCase struct {
	k3sRepository repositories.K3sRepository
}

func NewK3sUseCase(k3sRepository repositories.K3sRepository) K3sUseCase {
	return &k3sUseCase{k3sRepository: k3sRepository}
}

func (u *k3sUseCase) ListK3sNodes() ([]K3sNode, error) {
	return u.k3sRepository.ListK3sNodes()
}

func (u *k3sUseCase) ListRunningPods() ([]Pod, error) {
	return u.k3sRepository.ListRunningPods()
}

func (u *k3sUseCase) ListPodsInNamespace(namespace string) ([]Pod, error) {
	return u.k3sRepository.ListPodsInNamespace(namespace)
}

func (u *k3sUseCase) ListK3sServices() ([]Service, error) {
	return u.k3sRepository.ListK3sServices()
}

func (u *k3sUseCase) ListK3sDeployments() ([]Deployment, error) {
	return u.k3sRepository.ListK3sDeployments()
}

func (u *k3sUseCase) GetPodLogs(pod string) (string, error) {
	return u.k3sRepository.GetPodLogs(pod)
}

func (u *k3sUseCase) GetClusterEvents() ([]Event, error) {
	return u.k3sRepository.GetClusterEvents()
}

func (u *k3sUseCase) ApplyYAMLManifest(manifest string) error {
	return u.k3sRepository.ApplyYAMLManifest(manifest)
}

func (u *k3sUseCase) DeleteResource(resourceID string) error {
	return u.k3sRepository.DeleteResource(resourceID)
} 
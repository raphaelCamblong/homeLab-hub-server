package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type HealthStatusUseCase interface {
	GetStatus() (*entities.StatusEntity, error)
}

type healthUseCase struct {
	statusRepository repositories.StatusRepository
}

func NewStatusUseCase(statusRepository repositories.StatusRepository) HealthStatusUseCase {
	return &healthUseCase{statusRepository: statusRepository}
}

func (u *healthUseCase) GetStatus() (*entities.StatusEntity, error) {
	return u.statusRepository.GetStatus()
}

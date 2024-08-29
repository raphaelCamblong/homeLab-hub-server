package usecase

import (
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/cloud"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type CloudUseCase interface {
	GetVmsData() (*[]entities.VMEntity, error)
	GetHostData() (*entities.HostEntity, error)
}

type cloudUseCase struct {
	cloudRepository repositories.CloudRepository
}

func NewCloudUseCase(cloudRepository repositories.CloudRepository) CloudUseCase {
	return &cloudUseCase{cloudRepository: cloudRepository}
}

func (c cloudUseCase) GetVmsData() (*[]entities.VMEntity, error) {
	return c.cloudRepository.GetVMs()
}

func (c cloudUseCase) GetHostData() (*entities.HostEntity, error) {
	return c.cloudRepository.GetMainHost()
}

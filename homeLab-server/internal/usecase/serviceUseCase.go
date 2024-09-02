package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ServiceUseCase interface {
	GetAllService() (*[]entities.ServiceEntity, error)
	GetServiceById(string) (*entities.ServiceEntity, error)
}

type serviceUseCase struct {
	iloRepository repositories.ServiceRepository
}

func NewServiceUseCase(iloRepository repositories.ServiceRepository) ServiceUseCase {
	return &serviceUseCase{iloRepository: iloRepository}
}

func (u *serviceUseCase) GetAllService() (*[]entities.ServiceEntity, error) {
	return u.iloRepository.GetAllService()
}

func (u *serviceUseCase) GetServiceById(id string) (*entities.ServiceEntity, error) {
	return u.iloRepository.GetServiceById(id)
}

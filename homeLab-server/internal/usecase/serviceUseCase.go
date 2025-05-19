package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ServiceUseCase interface {
	ListServices() ([]entities.ServiceEntity, error)
	GetService(id string) (entities.ServiceEntity, error)
	CreateService(service entities.ServiceEntity) (entities.ServiceEntity, error)
	UpdateService(id string, service entities.ServiceEntity) (entities.ServiceEntity, error)
	DeleteService(id string) error
}

type serviceUseCase struct {
	serviceRepo repositories.ServiceRepository
}

func NewServiceUseCase(serviceRepo repositories.ServiceRepository) ServiceUseCase {
	return &serviceUseCase{serviceRepo: serviceRepo}
}

func (u *serviceUseCase) ListServices() ([]entities.ServiceEntity, error) {
	return u.serviceRepo.List()
}

func (u *serviceUseCase) GetService(id string) (entities.ServiceEntity, error) {
	return u.serviceRepo.Get(id)
}

func (u *serviceUseCase) CreateService(service entities.ServiceEntity) (entities.ServiceEntity, error) {
	return u.serviceRepo.CreateService(service)
}

func (u *serviceUseCase) UpdateService(id string, service entities.ServiceEntity) (entities.ServiceEntity, error) {
	return u.serviceRepo.UpdateService(id, service)
}

func (u *serviceUseCase) DeleteService(id string) error {
	return u.serviceRepo.DeleteService(id)
}

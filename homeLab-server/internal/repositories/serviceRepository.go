package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ServiceRepository interface {
	GetAllService() (*[]entities.ServiceEntity, error)
	GetServiceById(string) (*entities.ServiceEntity, error)
}

type serviceRepository struct {
	db database.Database
}

func NewServiceRepository(db database.Database) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetAllService() (*[]entities.ServiceEntity, error) {
	var services []entities.ServiceEntity
	if err := r.db.GetDb().Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}

func (r *serviceRepository) GetServiceById(id string) (*entities.ServiceEntity, error) {
	var service entities.ServiceEntity
	if err := r.db.GetDb().Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ServiceRepository interface {
	List() ([]entities.ServiceEntity, error)
	Get(id string) (entities.ServiceEntity, error)
	CreateService(service entities.ServiceEntity) (entities.ServiceEntity, error)
	UpdateService(id string, service entities.ServiceEntity) (entities.ServiceEntity, error)
	DeleteService(id string) error
}

type serviceRepository struct {
	db database.Database
}

func NewServiceRepository(db database.Database) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) List() ([]entities.ServiceEntity, error) {
	var services []entities.ServiceEntity
	result := r.db.Get().Find(&services)
	return services, result.Error
}

func (r *serviceRepository) Get(id string) (entities.ServiceEntity, error) {
	var service entities.ServiceEntity
	result := r.db.Get().First(&service, "id = ?", id)
	return service, result.Error
}

func (r *serviceRepository) CreateService(service entities.ServiceEntity) (entities.ServiceEntity, error) {
	result := r.db.Get().Create(&service)
	return service, result.Error
}

func (r *serviceRepository) UpdateService(id string, service entities.ServiceEntity) (entities.ServiceEntity, error) {
	var existingService entities.ServiceEntity
	if err := r.db.Get().First(&existingService, "id = ?", id).Error; err != nil {
		return entities.ServiceEntity{}, err
	}

	result := r.db.Get().Model(&existingService).Updates(service)
	if result.Error != nil {
		return entities.ServiceEntity{}, result.Error
	}

	return r.Get(id)
}

func (r *serviceRepository) DeleteService(id string) error {
	result := r.db.Get().Delete(&entities.ServiceEntity{}, "id = ?", id)
	return result.Error
}

package repositories

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities/ilo"
)

type ILORepository interface {
	GetThermal() (*entities.ThermalEntity, error)
}

type iloRepository struct {
	redfishRepository RedfishRepository
}

func NewThermalRepository(redfishRepository RedfishRepository) ILORepository {
	return &iloRepository{redfishRepository: redfishRepository}
}

func (r *iloRepository) GetThermal() (*entities.ThermalEntity, error) {
	if err := r.redfishRepository.UseSession(); err != nil {
		return nil, err
	}
	return r.redfishRepository.GetThermalData()
}

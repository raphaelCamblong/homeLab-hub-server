package repositories

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ThermalRepository interface {
	GetThermal(ctx *gin.Context) (*entities.ThermalEntity, error)
}

type thermalRepository struct {
	redfishRepository RedfishRepository
}

func NewThermalRepository(redfishRepository RedfishRepository) ThermalRepository {
	return &thermalRepository{redfishRepository: redfishRepository}
}

func (r *thermalRepository) GetThermal(_ *gin.Context) (*entities.ThermalEntity, error) {
	if !r.redfishRepository.IsSessionOpen() {
		if err := r.redfishRepository.CreateSession(); err != nil {
			return nil, err
		}
	}
	return r.redfishRepository.GetThermalData()
}

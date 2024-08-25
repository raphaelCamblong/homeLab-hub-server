package usecase

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ThermalUseCase interface {
	GetThermal(ctx *gin.Context) (*entities.ThermalEntity, error)
}

type thermalUseCase struct {
	thermalRepository repositories.ThermalRepository
}

func NewThermalUseCase(thermalRepository repositories.ThermalRepository) ThermalUseCase {
	return &thermalUseCase{thermalRepository: thermalRepository}
}

func (u *thermalUseCase) GetThermal(ctx *gin.Context) (*entities.ThermalEntity, error) {
	return u.thermalRepository.GetThermal(ctx)
}

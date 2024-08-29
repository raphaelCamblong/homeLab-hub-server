package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities/ilo"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type IlohUseCase interface {
	GetThermal() (*entities.ThermalEntity, error)
}

type iloUseCase struct {
	iloRepository repositories.ILORepository
}

func NewRewIloUseCase(iloRepository repositories.ILORepository) IlohUseCase {
	return &iloUseCase{iloRepository: iloRepository}
}

func (u *iloUseCase) GetThermal() (*entities.ThermalEntity, error) {
	return u.iloRepository.GetThermal()
}

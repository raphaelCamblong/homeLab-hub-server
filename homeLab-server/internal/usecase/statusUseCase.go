package usecase

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type StatusUseCase interface {
	GetStatus(ctx *gin.Context) (*entities.StatusEntity, error)
}

type statusUseCase struct {
	statusRepository repositories.StatusRepository
}

func NewStatusUseCase(statusRepository repositories.StatusRepository) StatusUseCase {
	return &statusUseCase{statusRepository: statusRepository}
}

func (u *statusUseCase) GetStatus(ctx *gin.Context) (*entities.StatusEntity, error) {
	return u.statusRepository.GetStatus(ctx)
}

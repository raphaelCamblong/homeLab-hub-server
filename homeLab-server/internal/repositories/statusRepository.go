package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type StatusRepository interface {
	GetStatus(ctx *gin.Context) (*entities.StatusEntity, error)
}

type statusRepository struct {
	database.Database
}

func NewStatusRepository(db database.Database) StatusRepository {
	return &statusRepository{db}
}

func (r *statusRepository) GetStatus(_ *gin.Context) (*entities.StatusEntity, error) {
	var status entities.StatusEntity
	result := r.GetDb().First(&status)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get status: %w", result.Error)
	}
	return &status, nil
}

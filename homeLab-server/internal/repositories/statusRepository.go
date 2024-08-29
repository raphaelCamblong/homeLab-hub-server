package repositories

import (
	"fmt"

	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type StatusRepository interface {
	GetStatus() (*entities.StatusEntity, error)
}

type statusRepository struct {
	database.Database
}

func NewStatusRepository(db database.Database) StatusRepository {
	return &statusRepository{db}
}

func (r *statusRepository) GetStatus() (*entities.StatusEntity, error) {
	var status entities.StatusEntity
	result := r.GetDb().First(&status)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get status: %w", result.Error)
	}
	return &status, nil
}

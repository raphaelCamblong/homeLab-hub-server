package repositories

import (
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type StatusRepository interface {
	GetStatus() (*entities.StatusEntity, error)
	GetHealthReport() (*[]entities.HealthEntity, error)
}

type statusRepository struct {
	*infrastructure.Infrastructure
}

func NewStatusRepository(infra *infrastructure.Infrastructure) StatusRepository {
	return &statusRepository{infra}
}

func (r *statusRepository) GetStatus() (*entities.StatusEntity, error) {
	var status entities.StatusEntity
	cfg := config.GetConfig()

	result := r.Db.GetDb().First(&status)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get status: %w", result.Error)
	}
	status.Version = cfg.App.Version
	return &status, nil
}

//func makeHealthReport(fetch func(option *externalHttpService.RequestOption) (*[]byte, error), ServiceName string ) entities.HealthEntity {
//	var report entities.HealthEntity
//	ctx := context.Background()
//
//
//	res, err := fetch(ctx)
//	if err != nil {
//
//	}
//	report.Name = ServiceName
//	report.IsOk =
//
//}

func (r *statusRepository) GetHealthReport() (entity *[]entities.HealthEntity, err error) {
	//	//httpInfra := r.ExternalHttpService
	//	//var healthReport []entities.HealthEntity
	//	//return healthReport, nil
	return nil, nil
}

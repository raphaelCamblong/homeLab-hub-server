package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type MonitoringUseCase interface {
	QueryPrometheus() (*PrometheusData, error)
	GetGrafanaDashboardData() (*GrafanaData, error)
	FetchActiveAlerts() ([]Alert, error)
}

type monitoringUseCase struct {
	monitoringRepository repositories.MonitoringRepository
}

func NewMonitoringUseCase(monitoringRepository repositories.MonitoringRepository) MonitoringUseCase {
	return &monitoringUseCase{monitoringRepository: monitoringRepository}
}

func (u *monitoringUseCase) QueryPrometheus() (*PrometheusData, error) {
	return u.monitoringRepository.QueryPrometheus()
}

func (u *monitoringUseCase) GetGrafanaDashboardData() (*GrafanaData, error) {
	return u.monitoringRepository.GetGrafanaDashboardData()
}

func (u *monitoringUseCase) FetchActiveAlerts() ([]Alert, error) {
	return u.monitoringRepository.FetchActiveAlerts()
} 
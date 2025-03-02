package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/externalHttpService"
	"github.com/prometheus/common/model"
)

type PrometheusRepository interface {
	GetMetrics(query string) (model.Value, error)
}

type prometheusRepository struct {
	Service externalHttpService.Prometheus
}

func NewPrometheusRepository(service externalHttpService.Prometheus) PrometheusRepository {
	return &prometheusRepository{Service: service}
}

func (r *prometheusRepository) GetMetrics(query string) (model.Value, error) {
	return r.Service.GetMetrics(query)
} 
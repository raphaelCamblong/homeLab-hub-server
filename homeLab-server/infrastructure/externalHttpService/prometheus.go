package externalHttpService

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/api/prometheus/v1/model"
)

type Prometheus interface {
	GetMetrics(query string) (model.Value, error)
}

type prometheusService struct {
	client v1.API
}

func NewPrometheusService(cfg NetworkConnection) (Prometheus, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.Host,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &prometheusService{client: v1.NewAPI(client)}, nil
}

func (p *prometheusService) GetMetrics(query string) (model.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := p.client.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus: %w", err)
	}
	return result, nil
} 
package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/internal/entities/ilo"
	"time"
)

type RedfishRepository interface {
	UseSession() error
	GetThermalData() (*ilo.ThermalEntity, error)
	GetPowerData() (*ilo.PowerEntity, error)
	GetPowerFastData() (*ilo.PowerEntity, error)
	GetHealthCheck() (*ilo.HealthEntity, error)
}

type redfishRepository struct {
	Cache   cache.Database
	Service client.Redfish
	ReqOpt  *client.RequestOption
}

func NewRedfishRepository(cache cache.Database, redfish client.Redfish) RedfishRepository {
	return &redfishRepository{
		Cache:   cache,
		Service: redfish,
	}
}

func (r *redfishRepository) UseSession() error {
	// Try to get token from the heap
	if r.ReqOpt != nil {
		return nil
	}
	// Try get token from redis cache
	opt, err := r.getCachedToken()
	if err == nil {
		r.ReqOpt = opt
		return nil
	}
	// Try get token from infra
	cfg := config.GetConfig()
	cred := client.Credentials{
		Username: *cfg.ExternalServicesCredential.Ilo.User,
		Password: *cfg.ExternalServicesCredential.Ilo.Key,
	}
	opt, err = r.Service.CreateSession(&cred)
	if err != nil {
		return err
	}
	r.ReqOpt = opt

	err = r.saveTokenCache(*opt)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (r *redfishRepository) getCachedToken() (*client.RequestOption, error) {
	ctx := context.Background()

	if r.Cache == nil {
		return nil, fmt.Errorf("can't get Redfish cached Token")
	}

	jsonData, err := r.Cache.GetClient().Get(ctx, "redfishToken").Result()
	if err != nil {
		return nil, fmt.Errorf("can't get Redfish cached Token")
	}
	var requestOption client.RequestOption
	err = json.Unmarshal([]byte(jsonData), &requestOption)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RequestOption from JSON: %w", err)
	}
	return &requestOption, nil
}

func (r *redfishRepository) saveTokenCache(option client.RequestOption) error {
	jsonData, err := json.Marshal(option)
	if err != nil {
		return fmt.Errorf("failed to marshal RequestOption: %w", err)
	}

	if r.Cache == nil {
		return fmt.Errorf("can't save Redfish Token")
	}
	ctx := context.Background()
	err = r.Cache.GetClient().Set(ctx, "redfishToken", jsonData, time.Minute*30).Err()
	if err != nil {
		return fmt.Errorf("can't save Redfish Token")
	}
	return nil
}

func (r *redfishRepository) GetThermalData() (*ilo.ThermalEntity, error) {
	thermalData, err := r.Service.GetThermalData(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return thermalData, nil
}

func (r *redfishRepository) GetPowerData() (*ilo.PowerEntity, error) {
	powerData, err := r.Service.GetPowerData(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power data: %w", err)
	}
	return powerData, nil
}

func (r *redfishRepository) GetPowerFastData() (*ilo.PowerEntity, error) {
	powerFastData, err := r.Service.GetPowerFastData(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power fast data: %w", err)
	}
	return powerFastData, nil
}

func (r *redfishRepository) GetHealthCheck() (*ilo.HealthEntity, error) {
	healthData, err := r.Service.GetHealthCheck(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve health data: %w", err)
	}
	return healthData, nil
}

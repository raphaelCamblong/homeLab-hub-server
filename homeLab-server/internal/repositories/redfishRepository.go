package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/externalHttpService"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"time"
)

type RedfishRepository interface {
	UseSession() error
	GetThermalData() (*entities.ThermalEntity, error)
	GetPowerData() (*entities.ThermalEntity, error)
}

type redfishRepository struct {
	Cache   cache.Database
	Service externalHttpService.Redfish
	ReqOpt  *externalHttpService.RequestOption
}

func NewRedfishRepository(cache cache.Database, redfish externalHttpService.Redfish) RedfishRepository {
	return &redfishRepository{
		Cache:   cache,
		Service: redfish,
	}
}

func (r *redfishRepository) UseSession() error {
	// Try get token from the heap
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
	cred := externalHttpService.Credentials{
		Username: cfg.ExternalServicesCredential.IloUsername,
		Password: cfg.ExternalServicesCredential.IloPassword,
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

func (r *redfishRepository) getCachedToken() (*externalHttpService.RequestOption, error) {
	ctx := context.Background()

	jsonData, err := r.Cache.GetClient().Get(ctx, "redfishToken").Result()
	if err != nil {
		return nil, fmt.Errorf("can't get Redfish cached Token")
	}
	var requestOption externalHttpService.RequestOption
	err = json.Unmarshal([]byte(jsonData), &requestOption)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RequestOption from JSON: %w", err)
	}

	return &requestOption, nil
}

func (r *redfishRepository) saveTokenCache(option externalHttpService.RequestOption) error {
	jsonData, err := json.Marshal(option)
	if err != nil {
		return fmt.Errorf("failed to marshal RequestOption: %w", err)
	}

	ctx := context.Background()
	err = r.Cache.GetClient().Set(ctx, "redfishToken", jsonData, time.Minute*30).Err()
	if err != nil {
		return fmt.Errorf("can't save Redfish Token")
	}
	return nil
}

func (r *redfishRepository) GetThermalData() (*entities.ThermalEntity, error) {
	bodyBytes, err := r.Service.GetThermalData(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("Failed to to retrieve thermal data: %w", err)
	}

	var thermalData entities.ThermalEntity
	thermalData, err = entities.UnmarshalThermalEntity(*bodyBytes)

	return &thermalData, err
}

func (r *redfishRepository) GetPowerData() (*entities.ThermalEntity, error) {
	_, err := r.Service.GetPowerData(r.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("Failed to to retrieve power data: %w", err)
	}

	return nil, err
}

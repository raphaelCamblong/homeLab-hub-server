package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/internal/entities/ilo"
	"time"
)

type ILORepository interface {
	GetThermal() (*entities.ThermalEntity, error)
	GetPower() (*entities.PowerEntity, error)
}

type iloRepository struct {
	redfishRepository RedfishRepository
	cache             cache.Database
}

func NewThermalRepository(redfishRepository RedfishRepository, cache cache.Database) ILORepository {
	return &iloRepository{redfishRepository: redfishRepository, cache: cache}
}

func (r *iloRepository) GetThermal() (*entities.ThermalEntity, error) {
	if err := r.redfishRepository.UseSession(); err != nil {
		return nil, err
	}
	return r.redfishRepository.GetThermalData()
}

func (r *iloRepository) GetPower() (*entities.PowerEntity, error) {
	powerCacheKey := "Ilo_power_data"
	ctx := context.Background()

	if r.cache != nil {

		data, err := r.cache.GetClient().Get(ctx, powerCacheKey).Result()
		if err == nil {

			powerEntity, err := entities.UnmarshalPowerEntity([]byte(data))
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal cached power data: %w", err)
			}
			return powerEntity, nil
		}
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
	}

	if err := r.redfishRepository.UseSession(); err != nil {
		return nil, err
	}
	powerEntity, err := r.redfishRepository.GetPowerData()
	if err != nil {
		return nil, err
	}

	marshaledData, err := powerEntity.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal power data: %w", err)
	}

	if r.cache == nil {
		return powerEntity, nil
	} else if err := r.cache.GetClient().Set(ctx, powerCacheKey, marshaledData, time.Minute*30).Err(); err != nil {
		return nil, fmt.Errorf("failed to cache power data: %w", err)
	}

	return powerEntity, nil
}

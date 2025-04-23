package repositories

import (
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/cloud"
	"strings"
)

type CloudRepository interface {
	GetVMs() (*[]entities.VMEntity, error)
	GetVM(string) (*entities.VMEntity, error)
	GetMainHost() (*entities.HostEntity, error)
}

type cloudRepository struct {
	xoRepository XenOrchestraRepository
}

func NewCloudRepository(xoRepository XenOrchestraRepository) CloudRepository {
	return &cloudRepository{xoRepository: xoRepository}
}

func (c *cloudRepository) GetVMs() (*[]entities.VMEntity, error) {
	c.xoRepository.UseSession()
	paths, err := c.xoRepository.GetAllVm()
	if err != nil {
		return nil, err
	}

	var vms []entities.VMEntity
	for _, path := range *paths {
		id := strings.Split(path, "/")[4]
		vm, err := c.xoRepository.GetVm(id)
		if err != nil {
			return nil, err
		}
		vms = append(vms, *vm)
	}
	return &vms, nil
}

func (c *cloudRepository) GetVM(id string) (*entities.VMEntity, error) {
	c.xoRepository.UseSession()
	vm, err := c.xoRepository.GetVm(id)
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (c *cloudRepository) GetMainHost() (*entities.HostEntity, error) {
	c.xoRepository.UseSession()
	paths, err := c.xoRepository.GetAllHost()
	if err != nil {
		return nil, err
	}

	for _, path := range *paths {
		id := strings.Split(path, "/")[4]
		host, err := c.xoRepository.GetHost(id)
		if err != nil {
			return nil, err
		}
		return host, nil
	}
	return nil, nil
}


package repositories

import (
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/cloud"
)

type XenOrchestraRepository interface {
	UseSession()
	GetAllVm() (*entities.XoRawPathEntity, error)
	GetVm(id string) (*entities.VMEntity, error)
	GetAllHost() (*entities.XoRawPathEntity, error)
	GetHost(id string) (*entities.HostEntity, error)
}

type xenOrchestraRepository struct {
	Cache   cache.Database
	Service client.XenOrchestra
	Db      database.Database
	ReqOpt  *client.RequestOption
}

func NewXenOrchestraRepository(cache cache.Database, xo client.XenOrchestra, db database.Database) XenOrchestraRepository {
	return &xenOrchestraRepository{
		Cache:   cache,
		Service: xo,
		Db:      db,
	}
}

func (x *xenOrchestraRepository) UseSession() {
	if x.ReqOpt != nil {
		return
	}

	x.ReqOpt = &client.RequestOption{
		AuthToken: *config.GetConfig().ExternalServicesCredential.XO.Key,
	}
}

func (x *xenOrchestraRepository) GetAllVm() (*entities.XoRawPathEntity, error) {
	bodyBytes, err := x.Service.GetAllVm(x.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to to retrieve Vms data: %w", err)
	}

	xoPath, err := entities.UnmarshalXoRawPathEntity(*bodyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to to unmarshal Vms data: %w", err)
	}
	return xoPath, err
}

func (x *xenOrchestraRepository) GetVm(id string) (*entities.VMEntity, error) {
	bodyBytes, err := x.Service.GetVm(id, x.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to to retrieve Vm: '%s' data: %w", id, err)
	}

	vm, err := entities.UnmarshalVMEntity(*bodyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to to unmarshal Host data: %w", err)
	}
	return vm, err
}

func (x *xenOrchestraRepository) GetAllHost() (*entities.XoRawPathEntity, error) {
	bodyBytes, err := x.Service.GetAllHost(x.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to to retrieve Hosts data: %w", err)
	}

	var xoPath *entities.XoRawPathEntity
	xoPath, err = entities.UnmarshalXoRawPathEntity(*bodyBytes)

	return xoPath, err
}

func (x *xenOrchestraRepository) GetHost(id string) (*entities.HostEntity, error) {
	bodyBytes, err := x.Service.GetHost(id, x.ReqOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to to retrieve host: '%s' data: %w", id, err)
	}

	var host *entities.HostEntity
	host, err = entities.UnmarshalHostEntity(*bodyBytes)

	return host, err
}



// ILO
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

func NewIloRepository(redfishRepository RedfishRepository, cache cache.Database) ILORepository {
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

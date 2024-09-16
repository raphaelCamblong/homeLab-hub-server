package repositories

import (
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/externalHttpService"
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
	Service externalHttpService.XenOrchestra
	Db      database.Database
	ReqOpt  *externalHttpService.RequestOption
}

func NewXenOrchestraRepository(cache cache.Database, xo externalHttpService.XenOrchestra, db database.Database) XenOrchestraRepository {
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

	x.ReqOpt = &externalHttpService.RequestOption{
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

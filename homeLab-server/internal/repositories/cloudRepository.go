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

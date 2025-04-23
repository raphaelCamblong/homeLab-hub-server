package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type Repositories struct {
	XenOrchestra repositories.XenOrchestraRepository
	Status       repositories.StatusRepository
	Ilo          repositories.ILORepository
	Redfish      repositories.RedfishRepository
	Cloud        repositories.CloudRepository
	Auth         repositories.AuthenticationRepository
	Service      repositories.ServiceRepository
}

func NewRepositories(infra *infrastructure.Infrastructure) *Repositories {
	repo := Repositories{
		XenOrchestra: repositories.NewXenOrchestraRepository(infra),
		Status:       repositories.NewStatusRepository(infra),
		Redfish:      repositories.NewRedfishRepository(infra),
		Auth:         repositories.NewAuthenticationRepository(infra),
		Service:      repositories.NewServiceRepository(infra),
	}

	repo.Ilo = repositories.NewIloRepository(infra)
	repo.Cloud = repositories.NewCloudRepository(infra)
	return &repo
}

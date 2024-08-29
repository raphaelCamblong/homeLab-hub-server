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
}

func NewRepositories(infra *infrastructure.Infrastructure) *Repositories {
	repo := Repositories{
		XenOrchestra: repositories.NewXenOrchestraRepository(infra.Cache, infra.ExternalHttpService.GetXenOrchestra(), infra.Db),
		Status:       repositories.NewStatusRepository(infra.Db),
		Redfish:      repositories.NewRedfishRepository(infra.Cache, infra.ExternalHttpService.GetRedfish()),
		Auth:         repositories.NewAuthenticationRepository(infra.Db),
	}
	repo.Ilo = repositories.NewThermalRepository(repo.Redfish, infra.Cache)
	repo.Cloud = repositories.NewCloudRepository(repo.XenOrchestra)
	return &repo
}

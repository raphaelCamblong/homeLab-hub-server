package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
)

type Repositories struct {
	User       UserRepository
	Service    ServiceRepository
	Pipeline   PipelineRepository
	Infra      InfrastructureRepository
	Storage    StorageRepository
	Network    NetworkRepository
	Cluster    ClusterRepository
	Monitoring MonitoringRepository
}

func NewRepositories(infra *infrastructure.Infrastructure) *Repositories {
	repo := Repositories{
		User:       NewUserRepository(infra.Db),
		Service:    NewServiceRepository(infra.Db),
		Pipeline:   NewPipelineRepository(infra.Db, infra.StreamHub, infra.Cron),
		Infra:      NewInfrastructureRepository(infra.Db),
		Storage:    NewStorageRepository(),
		Network:    NewNetworkRepository(),
		Cluster:    NewClusterRepository(),
		Monitoring: NewMonitoringRepository(),
	}
	return &repo
}

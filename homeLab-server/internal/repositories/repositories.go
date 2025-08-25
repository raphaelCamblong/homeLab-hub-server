package repositories

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
)

type Repositories struct {
	User         UserRepository
	Service      ServiceRepository
	Pipeline     PipelineRepository
	Infra        InfrastructureRepository
	Storage      StorageRepository
	Cluster      ClusterRepository
	Monitoring   MonitoringRepository
	Notification NotificationRepository
}

func NewRepositories(infra *infrastructure.Infrastructure) *Repositories {
	repo := Repositories{
		User:         NewUserRepository(infra.Db),
		Service:      NewServiceRepository(infra.Db),
		Pipeline:     NewPipelineRepository(infra.Db, infra.StreamHub, infra.Cron),
		Infra:        NewInfrastructureRepository(infra.Db),
		Storage:      NewStorageRepository(),
		Cluster:      NewClusterRepository(infra.ExternalHttpService.GetK8sClient()),
		Monitoring:   NewMonitoringRepository(),
		Notification: NewNotificationRepository(infra.Db, infra.StreamHub),
	}
	return &repo
}

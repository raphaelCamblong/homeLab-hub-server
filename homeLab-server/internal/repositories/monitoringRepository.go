package repositories

type MonitoringRepository interface{}

type monitoringRepository struct{}

func NewMonitoringRepository() MonitoringRepository {
	return &monitoringRepository{}
}

package repositories

import (
	"time"

	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type InfrastructureRepository interface {
	// Basic CRUD operations
	List() ([]entities.InfrastructureEntity, error)
	Get(id string) (entities.InfrastructureEntity, error)
	Create(infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error)
	Update(id string, infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error)
	Delete(id string) error

	// Infrastructure specific operations
	GetByHostname(hostname string) (entities.InfrastructureEntity, error)
	GetByIPAddress(ipAddress string) (entities.InfrastructureEntity, error)
	GetByMACAddress(macAddress string) (entities.InfrastructureEntity, error)
	ListByProvider(provider string) ([]entities.InfrastructureEntity, error)
	ListByType(infraType string) ([]entities.InfrastructureEntity, error)
	ListByRegion(region string) ([]entities.InfrastructureEntity, error)
	UpdateHealthStatus(id string, status string, latencyMS int) error
	UpdateFeatures(id string, features map[string]interface{}) error

	// Background job specific operations
	UpdateObservabilityMetrics(id string, uptime float64, latencyMS int, healthStatus string) error
	UpdateResourceMetrics(id string, cpuCores int, ramMB int, diskSizeGB int) error
	BulkUpdateHealthStatus(ids []string, status string) error
	UpsertByHostname(infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error)
}

type infrastructureRepository struct {
	db database.Database
}

func NewInfrastructureRepository(db database.Database) InfrastructureRepository {
	return &infrastructureRepository{db: db}
}

func (r *infrastructureRepository) List() ([]entities.InfrastructureEntity, error) {
	var infras []entities.InfrastructureEntity
	result := r.db.Get().Find(&infras)
	return infras, result.Error
}

func (r *infrastructureRepository) Get(id string) (entities.InfrastructureEntity, error) {
	var infra entities.InfrastructureEntity
	result := r.db.Get().First(&infra, "id = ?", id)
	return infra, result.Error
}

func (r *infrastructureRepository) Create(infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error) {
	result := r.db.Get().Create(&infra)
	return infra, result.Error
}

func (r *infrastructureRepository) Update(id string, infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error) {
	var existingInfra entities.InfrastructureEntity
	if err := r.db.Get().First(&existingInfra, "id = ?", id).Error; err != nil {
		return entities.InfrastructureEntity{}, err
	}

	result := r.db.Get().Model(&existingInfra).Updates(infra)
	if result.Error != nil {
		return entities.InfrastructureEntity{}, result.Error
	}

	return r.Get(id)
}

func (r *infrastructureRepository) Delete(id string) error {
	result := r.db.Get().Delete(&entities.InfrastructureEntity{}, "id = ?", id)
	return result.Error
}

func (r *infrastructureRepository) GetByHostname(hostname string) (entities.InfrastructureEntity, error) {
	var infra entities.InfrastructureEntity
	result := r.db.Get().First(&infra, "hostname = ?", hostname)
	return infra, result.Error
}

func (r *infrastructureRepository) GetByIPAddress(ipAddress string) (entities.InfrastructureEntity, error) {
	var infra entities.InfrastructureEntity
	result := r.db.Get().First(&infra, "ip_address = ?", ipAddress)
	return infra, result.Error
}

func (r *infrastructureRepository) GetByMACAddress(macAddress string) (entities.InfrastructureEntity, error) {
	var infra entities.InfrastructureEntity
	result := r.db.Get().First(&infra, "mac_address = ?", macAddress)
	return infra, result.Error
}

func (r *infrastructureRepository) ListByProvider(provider string) ([]entities.InfrastructureEntity, error) {
	var infras []entities.InfrastructureEntity
	result := r.db.Get().Where("provider = ?", provider).Find(&infras)
	return infras, result.Error
}

func (r *infrastructureRepository) ListByType(infraType string) ([]entities.InfrastructureEntity, error) {
	var infras []entities.InfrastructureEntity
	result := r.db.Get().Where("type = ?", infraType).Find(&infras)
	return infras, result.Error
}

func (r *infrastructureRepository) ListByRegion(region string) ([]entities.InfrastructureEntity, error) {
	var infras []entities.InfrastructureEntity
	result := r.db.Get().Where("region = ?", region).Find(&infras)
	return infras, result.Error
}

func (r *infrastructureRepository) UpdateHealthStatus(id string, status string, latencyMS int) error {
	now := time.Now()
	updates := map[string]interface{}{
		"observability": entities.InfrastructureObservability{
			LastCheckedAt: &now,
			HealthStatus:  &status,
			LatencyMS:     &latencyMS,
		},
	}

	result := r.db.Get().Model(&entities.InfrastructureEntity{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (r *infrastructureRepository) UpdateFeatures(id string, features map[string]interface{}) error {
	result := r.db.Get().Model(&entities.InfrastructureEntity{}).
		Where("id = ?", id).
		Update("features", features)
	return result.Error
}

func (r *infrastructureRepository) UpdateObservabilityMetrics(id string, uptime float64, latencyMS int, healthStatus string) error {
	now := time.Now()
	updates := entities.InfrastructureObservability{
		LastCheckedAt: &now,
		UptimePercent: &uptime,
		LatencyMS:     &latencyMS,
		HealthStatus:  &healthStatus,
	}

	result := r.db.Get().Model(&entities.InfrastructureEntity{}).
		Where("id = ?", id).
		Update("observability", updates)
	return result.Error
}

func (r *infrastructureRepository) UpdateResourceMetrics(id string, cpuCores int, ramMB int, diskSizeGB int) error {
	updates := map[string]interface{}{
		"cpu_cores":    cpuCores,
		"ram_mb":       ramMB,
		"disk_size_gb": diskSizeGB,
	}

	result := r.db.Get().Model(&entities.InfrastructureEntity{}).
		Where("id = ?", id).
		Updates(updates)
	return result.Error
}

func (r *infrastructureRepository) BulkUpdateHealthStatus(ids []string, status string) error {
	now := time.Now()
	updates := entities.InfrastructureObservability{
		LastCheckedAt: &now,
		HealthStatus:  &status,
	}

	result := r.db.Get().Model(&entities.InfrastructureEntity{}).
		Where("id IN ?", ids).
		Update("observability", updates)
	return result.Error
}

func (r *infrastructureRepository) UpsertByHostname(infra entities.InfrastructureEntity) (entities.InfrastructureEntity, error) {
	var existing entities.InfrastructureEntity
	result := r.db.Get().Where("hostname = ?", infra.Hostname).First(&existing)

	if result.Error != nil {
		return r.Create(infra)
	}

	updateResult := r.db.Get().Model(&existing).Updates(infra)
	if updateResult.Error != nil {
		return entities.InfrastructureEntity{}, updateResult.Error
	}

	return r.Get(string(existing.ID))
}

package repositories

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type NetworkRepository interface {
	ListDevices() ([]entities.Device, error)
	ListNetworkInterfaces() ([]entities.NetworkInterface, error)
	GetInterfaceDetails(id string) (*entities.NetworkInterfaceDetails, error)
	ViewFirewallRules() ([]entities.FirewallRule, error)
	ModifyFirewallRules(rules []entities.FirewallRule) error
	ViewDNSSettings() (*entities.DNSSettings, error)
	RunNetworkSpeedTest() (*entities.SpeedTestResult, error)
}

type networkRepository struct {
	// Add any necessary fields, such as database connections
}

func NewNetworkRepository() NetworkRepository {
	return &networkRepository{}
}

func (r *networkRepository) ListDevices() ([]entities.Device, error) {
	// Implement the logic to list devices
}

func (r *networkRepository) ListNetworkInterfaces() ([]entities.NetworkInterface, error) {
	// Implement the logic to list network interfaces
}

func (r *networkRepository) GetInterfaceDetails(id string) (*entities.NetworkInterfaceDetails, error) {
	// Implement the logic to get interface details
}

func (r *networkRepository) ViewFirewallRules() ([]entities.FirewallRule, error) {
	// Implement the logic to view firewall rules
}

func (r *networkRepository) ModifyFirewallRules(rules []entities.FirewallRule) error {
	// Implement the logic to modify firewall rules
}

func (r *networkRepository) ViewDNSSettings() (*entities.DNSSettings, error) {
	// Implement the logic to view DNS settings
}


func (r *serviceRepository) GetAllService() (*[]entities.ServiceEntity, error) {
	var services []entities.ServiceEntity
	if err := r.db.GetDb().Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}

func (r *serviceRepository) GetServiceById(id string) (*entities.ServiceEntity, error) {
	var service entities.ServiceEntity
	if err := r.db.GetDb().Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}
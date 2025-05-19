package entities

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func UnmarshalInfrastructureEntity(data []byte) (InfrastructureEntity, error) {
	var r InfrastructureEntity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *InfrastructureEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type InfrastructureEntity struct {
	gorm.Model

	// Identity (core required fields)
	Name     string `json:"name" gorm:"not null"`
	Hostname string `json:"hostname" gorm:"not null"`

	// Identity (optional fields)
	Description *string `json:"description"`
	IPAddress   *string `json:"ip_address"`
	MACAddress  *string `json:"mac_address"`

	// Platform Metadata
	Provider *string           `json:"provider"`
	Type     *string           `json:"type"`
	Region   *string           `json:"region"`
	OS       *string           `json:"os"`
	Arch     *string           `json:"arch"`
	Tags     []string          `json:"tags" gorm:"type:json"`
	Labels   map[string]string `json:"labels" gorm:"type:json"`

	// Specs
	CPUModel     *string `json:"cpu_model"`
	CPUCores     *int    `json:"cpu_cores"`
	RAMMB        *int    `json:"ram_mb"`
	DiskSizeGB   *int    `json:"disk_size_gb"`
	GPUCount     *int    `json:"gpu_count"`
	NetworkSpeed *string `json:"network_speed"`
	Interfaces   *int    `json:"interfaces"`

	// Additional Specs
	StorageType        *string   `json:"storage_type"`        // SSD, HDD, NVMe
	VirtualizationType *string   `json:"virtualization_type"` // KVM, VMware, Docker
	PowerConsumptionW  *int      `json:"power_consumption_w"`
	RackLocation       *string   `json:"rack_location"`
	PhysicalLocation   *Location `json:"physical_location" gorm:"embedded"`

	// Network Configuration
	NetworkConfig *NetworkConfig `json:"network_config" gorm:"embedded"`
	DNSServers    []string       `json:"dns_servers" gorm:"type:json"`
	SubnetMask    *string        `json:"subnet_mask"`
	Gateway       *string        `json:"gateway"`
	VLANs         []string       `json:"vlans" gorm:"type:json"`

	// Security
	SecurityConfig *SecurityConfig `json:"security_config" gorm:"embedded"`

	// Observability
	Observability *InfrastructureObservability `gorm:"embedded" json:"observability"`

	// Dynamic Frontend Features
	Features map[string]interface{} `gorm:"type:json" json:"features"`
}

type Location struct {
	Building    *string `json:"building"`
	Floor       *string `json:"floor"`
	Room        *string `json:"room"`
	Coordinates *string `json:"coordinates"`
}

type NetworkConfig struct {
	PrimaryInterface   *string  `json:"primary_interface"`
	BondingMode        *string  `json:"bonding_mode"`
	NetworkZone        *string  `json:"network_zone"`
	FirewallRules      []string `json:"firewall_rules" gorm:"type:json"`
	LoadBalanced       *bool    `json:"load_balanced"`
	ProxyEnabled       *bool    `json:"proxy_enabled"`
	BandwidthLimitMbps *int     `json:"bandwidth_limit_mbps"`
}

type SecurityConfig struct {
	SSHEnabled      *bool      `json:"ssh_enabled"`
	SSHPort         *int       `json:"ssh_port"`
	FirewallEnabled *bool      `json:"firewall_enabled"`
	TPMEnabled      *bool      `json:"tpm_enabled"`
	SELinuxMode     *string    `json:"selinux_mode"`
	LastPatchedAt   *time.Time `json:"last_patched_at"`
	SecurityLevel   *string    `json:"security_level"`
	Certificates    []string   `json:"certificates" gorm:"type:json"`
}

type InfrastructureObservability struct {
	LastCheckedAt *time.Time `json:"last_checked_at"`
	UptimePercent *float64   `json:"uptime_percent"`
	LatencyMS     *int       `json:"latency_ms"`
	HealthStatus  *string    `json:"health_status"`

	// Additional monitoring metrics
	CPUUsagePercent       *float64 `json:"cpu_usage_percent"`
	MemoryUsagePercent    *float64 `json:"memory_usage_percent"`
	DiskUsagePercent      *float64 `json:"disk_usage_percent"`
	NetworkInBytesPerSec  *float64 `json:"network_in_bytes_per_sec"`
	NetworkOutBytesPerSec *float64 `json:"network_out_bytes_per_sec"`
	IOPSRead              *int     `json:"iops_read"`
	IOPSWrite             *int     `json:"iops_write"`
	Temperature           *float64 `json:"temperature"`
}

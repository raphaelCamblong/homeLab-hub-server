package entities

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func UnmarshalServiceEntity(data []byte) (ServiceEntity, error) {
	var r ServiceEntity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ServiceEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Main Entity
type ServiceEntity struct {
	gorm.Model

	// Required core fields
	Name      string `json:"name" gorm:"not null"`
	Namespace string `json:"namespace" gorm:"not null"`
	Status    string `json:"status" gorm:"not null"`
	Type      string `json:"type" gorm:"not null"`
	IP        string `json:"ip" gorm:"not null"`
	Port      int64  `json:"port" gorm:"not null"`

	// Optional
	Tags        string `json:"tags"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`

	// Subcategories
	Metadata      Metadata `gorm:"embedded" json:"metadata"`
	Observability `gorm:"embedded" json:"observability"`
	Security      Security    `gorm:"embedded" json:"security"`
	APIInfo       APIInfo     `gorm:"embedded" json:"api_info"`
	NetworkInfo   NetworkInfo `gorm:"embedded" json:"network_info"`
}

type NetworkInfo struct {
	Domain string `json:"domain"`
	Port   int64  `json:"port"`
	IP     string `json:"ip"`
	URL    string `json:"url"`
}

// Metadata Info
type Metadata struct {
	Owner       string `json:"owner"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// Observability Info
type Observability struct {
	DeployedAt    *time.Time `json:"deployed_at"`
	LastCheckedAt *time.Time `json:"last_checked_at"`
	UptimePercent float64    `json:"uptime_percent"`
	LatencyMS     int        `json:"latency_ms"`
	HealthStatus  string     `json:"health_status"`
}

// Security Info
type Security struct {
	AuthRequired bool `json:"auth_required"`
	TLS          bool `json:"tls"`
}

// API Info
type APIInfo struct {
	BaseURL    string `json:"base_url"`
	OpenAPIURL string `json:"openapi_url"`
}

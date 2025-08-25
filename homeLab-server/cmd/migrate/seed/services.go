package seed

import (
	"time"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

func SeedServices(db database.Database) {
	if db.Get().Migrator().HasTable(&entities.ServiceEntity{}) {
		logrus.Info("✅ Services table already exists")
		return
	}

	logrus.Info("Seeding services...")

	now := time.Now()
	services := []entities.ServiceEntity{
		{
			Name:        "Prometheus",
			Namespace:   "monitoring",
			Status:      "running",
			Type:        "monitoring",
			IP:          "192.168.1.100",
			Port:        9090,
			Tags:        "metrics,monitoring,observability",
			Description: "Prometheus monitoring system & time series database",
			LogoURL:     "https://prometheus.io/assets/prometheus_logo.png",
			NetworkInfo: entities.NetworkInfo{
				Domain: "prometheus.homelab.com",
				Port:   9090,
				IP:     "192.168.1.100",
				URL:    "http://prometheus.homelab.com:9090",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "2.45.0",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.99,
				LatencyMS:     5,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "http://prometheus:9090/api/v1",
				OpenAPIURL: "http://prometheus:9090/api/v1/openapi.json",
			},
		},
		{
			Name:        "MinIO",
			Namespace:   "storage",
			Status:      "running",
			Type:        "object-storage",
			IP:          "192.168.1.101",
			Port:        9000,
			Tags:        "storage,s3,backup",
			Description: "High Performance Object Storage",
			LogoURL:     "https://min.io/resources/img/logo.svg",
			NetworkInfo: entities.NetworkInfo{
				Domain: "minio.homelab.com",
				Port:   9000,
				IP:     "192.168.1.101",
				URL:    "http://minio.homelab.com:9000",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "RELEASE.2024-03-15",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.95,
				LatencyMS:     10,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "http://minio:9000",
				OpenAPIURL: "",
			},
		},
		{
			Name:        "Keycloak",
			Namespace:   "auth",
			Status:      "running",
			Type:        "authentication",
			IP:          "192.168.1.102",
			Port:        8080,
			Tags:        "auth,sso,identity",
			Description: "Open Source Identity and Access Management",
			LogoURL:     "https://www.keycloak.org/resources/images/keycloak_logo.png",
			NetworkInfo: entities.NetworkInfo{
				Domain: "keycloak.homelab.com",
				Port:   8080,
				IP:     "192.168.1.102",
				URL:    "http://keycloak.homelab.com:8080",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "22.0.1",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.98,
				LatencyMS:     15,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "http://keycloak:8080/auth",
				OpenAPIURL: "http://keycloak:8080/auth/realms/master/.well-known/openid-configuration",
			},
		},
		{
			Name:        "Grafana",
			Namespace:   "monitoring",
			Status:      "running",
			Type:        "visualization",
			IP:          "192.168.1.103",
			Port:        3000,
			Tags:        "dashboard,monitoring,visualization",
			Description: "Analytics and interactive visualization web application",
			LogoURL:     "https://grafana.com/static/img/logos/grafana_logo.svg",
			NetworkInfo: entities.NetworkInfo{
				Domain: "grafana.homelab.com",
				Port:   3000,
				IP:     "192.168.1.103",
				URL:    "http://grafana.homelab.com:3000",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "10.2.3",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.97,
				LatencyMS:     8,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "http://grafana:3000/api",
				OpenAPIURL: "",
			},
		},
		{
			Name:        "PostgreSQL",
			Namespace:   "database",
			Status:      "running",
			Type:        "database",
			IP:          "192.168.1.104",
			Port:        5432,
			Tags:        "database,sql,persistence",
			Description: "Advanced open-source database",
			LogoURL:     "https://www.postgresql.org/media/img/about/press/elephant.png",
			NetworkInfo: entities.NetworkInfo{
				Domain: "postgres.homelab.com",
				Port:   5432,
				IP:     "192.168.1.104",
				URL:    "http://postgres.homelab.com:5432",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "16.2",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.999,
				LatencyMS:     3,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "",
				OpenAPIURL: "",
			},
		},
		{
			Name:        "Traefik",
			Namespace:   "networking",
			Status:      "running",
			Type:        "proxy",
			IP:          "192.168.1.105",
			Port:        80,
			Tags:        "proxy,loadbalancer,edge-router",
			Description: "Cloud Native Edge Router",
			LogoURL:     "https://traefik.io/static/traefik-proxy-logo.svg",
			NetworkInfo: entities.NetworkInfo{
				Domain: "traefik.homelab.com",
				Port:   80,
				IP:     "192.168.1.105",
				URL:    "http://traefik.homelab.com:80",
			},
			Metadata: entities.Metadata{
				Owner:       "admin",
				Version:     "2.10.5",
				Environment: "production",
			},
			Observability: entities.Observability{
				DeployedAt:    &now,
				LastCheckedAt: &now,
				UptimePercent: 99.99,
				LatencyMS:     2,
				HealthStatus:  "healthy",
			},
			Security: entities.Security{
				AuthRequired: true,
				TLS:          true,
			},
			APIInfo: entities.APIInfo{
				BaseURL:    "http://traefik:8080/api",
				OpenAPIURL: "http://traefik:8080/api/rawdata",
			},
		},
	}

	for _, service := range services {
		if err := db.Get().Create(&service).Error; err != nil {
			logrus.Errorf("❌ Failed to create service %s: %v", service.Name, err)
		} else {
			logrus.Infof("✅ Created service: %s", service.Name)
		}
	}
}

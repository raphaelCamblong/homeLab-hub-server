package seed

import (
	"encoding/json"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ConfigFile struct {
	Steps     []StepConfig     `json:"steps"`
	Pipelines []PipelineConfig `json:"pipelines"`
}

type StepConfig struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

type PipelineConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

func LoadDataFromJson(db database.Database) {
	jsonFile, err := os.Open("pipeline.config.json")
	if err != nil {
		logrus.Errorf("❌ Failed to open pipeline.config.json: %v", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		logrus.Errorf("❌ Failed to read pipeline.config.json: %v", err)
		return
	}

	var config ConfigFile
	if err := json.Unmarshal(byteValue, &config); err != nil {
		logrus.Errorf("❌ Failed to parse pipeline.config.json: %v", err)
		return
	}

	// Create step templates from the JSON data
	for _, stepConfig := range config.Steps {
		// Check if step already exists
		var existingStep entities.StepTemplate
		result := db.Get().Where("name = ?", stepConfig.Name).First(&existingStep)
		if result.Error == nil {
			logrus.Infof("⏩ Step %s already exists, skipping creation", stepConfig.Name)
			continue
		}

		configJSON, err := json.Marshal(stepConfig.Config)
		if err != nil {
			logrus.Errorf("❌ Failed to marshal step config for %s: %v", stepConfig.Name, err)
			continue
		}

		step := entities.StepTemplate{
			Type:        stepConfig.Type,
			Name:        stepConfig.Name,
			Description: stepConfig.Description,
			Config:      string(configJSON),
		}

		if err := db.Get().Create(&step).Error; err != nil {
			logrus.Errorf("❌ Failed to create step %s: %v", step.Name, err)
			continue
		}
		logrus.Infof("✅ Created step template: %s", step.Name)
	}

	// Create pipeline templates from the JSON data
	for _, pipelineConfig := range config.Pipelines {
		// Check if pipeline already exists
		var existingPipeline entities.PipelineTemplate
		result := db.Get().Where("name = ?", pipelineConfig.Name).First(&existingPipeline)
		if result.Error == nil {
			logrus.Infof("⏩ Pipeline %s already exists, skipping creation", pipelineConfig.Name)
			continue
		}

		pipeline := entities.PipelineTemplate{
			Name:        pipelineConfig.Name,
			Description: pipelineConfig.Description,
		}

		if err := db.Get().Create(&pipeline).Error; err != nil {
			logrus.Errorf("❌ Failed to create pipeline template %s: %v", pipeline.Name, err)
			continue
		}

		// Associate steps with the pipeline
		for _, stepName := range pipelineConfig.Steps {
			var step entities.StepTemplate
			if err := db.Get().Where("name = ?", stepName).First(&step).Error; err != nil {
				logrus.Errorf("❌ Failed to find step %s for pipeline %s: %v", stepName, pipeline.Name, err)
				continue
			}

			// Check if association already exists
			var steps []entities.StepTemplate
			if err := db.Get().Model(&pipeline).Association("Steps").Find(&steps); err != nil {
				logrus.Errorf("❌ Failed to get existing steps for pipeline %s: %v", pipeline.Name, err)
				continue
			}

			exists := false
			for _, s := range steps {
				if s.ID == step.ID {
					exists = true
					break
				}
			}
			if exists {
				logrus.Infof("⏩ Step %s already associated with pipeline %s, skipping", stepName, pipeline.Name)
				continue
			}

			if err := db.Get().Model(&pipeline).Association("Steps").Append(&step); err != nil {
				logrus.Errorf("❌ Failed to associate step %s with pipeline %s: %v", stepName, pipeline.Name, err)
				continue
			}
			logrus.Infof("✅ Associated step %s with pipeline %s", stepName, pipeline.Name)
		}

		logrus.Infof("✅ Created pipeline template: %s with %d steps", pipeline.Name, len(pipelineConfig.Steps))
	}
}

func SeedPipelines(db database.Database) {
	LoadDataFromJson(db)
	// First, create step templates that can be reused
	// commonSteps := []entities.StepTemplate{
	// 	{Type: "lambda", Name: "Check disk space", Description: "Verify available disk space", Config: `{"min_space_gb": 10}`},
	// 	{Type: "lambda", Name: "Initialize", Description: "Initialize the process", Config: `{"timeout_seconds": 30}`},
	// 	{Type: "lambda", Name: "Generate report", Description: "Create execution report", Config: `{"format": "json"}`},
	// 	{Type: "lambda", Name: "Clean up", Description: "Clean up temporary files and resources", Config: `{"delete_temp": true}`},
	// }

	// for _, step := range commonSteps {
	// 	if err := db.Get().Create(&step).Error; err != nil {
	// 		logrus.Errorf("❌ Failed to create common step %s: %v", step.Name, err)
	// 	}
	// }

	// // Create pipeline-specific step templates
	// pipelineSteps := map[string][]entities.StepTemplate{
	// 	"System Update": {
	// 		{Name: "Check for updates", Description: "Check for available system updates", Config: `{"update_type": "security"}`},
	// 		{Name: "Download packages", Description: "Download update packages", Config: `{"concurrent_downloads": 3}`},
	// 		{Name: "Install updates", Description: "Install downloaded updates", Config: `{"auto_restart": false}`},
	// 	},
	// 	"Database Backup": {
	// 		{Name: "Stop services", Description: "Stop database services", Config: `{"services": ["postgresql", "mysql"]}`},
	// 		{Name: "Create backup", Description: "Create database backup files", Config: `{"compression": true}`},
	// 		{Name: "Upload to storage", Description: "Upload backup to remote storage", Config: `{"storage_type": "s3"}`},
	// 		{Name: "Restart services", Description: "Restart database services", Config: `{"timeout_seconds": 60}`},
	// 	},
	// 	"Network Security Scan": {
	// 		{Name: "Initialize scanner", Description: "Initialize security scanner", Config: `{"scanner": "nmap"}`},
	// 		{Name: "Port scan", Description: "Scan network ports", Config: `{"port_range": "1-65535"}`},
	// 		{Name: "Vulnerability check", Description: "Check for vulnerabilities", Config: `{"severity": "high"}`},
	// 	},
	// 	"Docker Container Cleanup": {
	// 		{Name: "List unused resources", Description: "List unused Docker resources", Config: `{"resource_types": ["containers", "images", "volumes"]}`},
	// 		{Name: "Stop unused containers", Description: "Stop inactive containers", Config: `{"timeout_seconds": 30}`},
	// 		{Name: "Remove containers", Description: "Remove stopped containers", Config: `{"force": false}`},
	// 		{Name: "Remove unused images", Description: "Remove unused Docker images", Config: `{"keep_tags": ["latest"]}`},
	// 		{Name: "Clean volumes", Description: "Clean unused Docker volumes", Config: `{"preserve_named": true}`},
	// 	},
	// 	"SSL Certificate Renewal": {
	// 		{Name: "Check certificate expiry", Description: "Check SSL certificate expiration", Config: `{"warn_days": 30}`},
	// 		{Name: "Generate new certificates", Description: "Generate new SSL certificates", Config: `{"provider": "letsencrypt"}`},
	// 		{Name: "Backup old certificates", Description: "Backup existing certificates", Config: `{"backup_location": "/etc/ssl/backup"}`},
	// 		{Name: "Install new certificates", Description: "Install new SSL certificates", Config: `{"restart_services": true}`},
	// 		{Name: "Reload web server", Description: "Reload web server configuration", Config: `{"graceful": true}`},
	// 	},
	// 	"System Health Check": {
	// 		{Name: "Check CPU usage", Description: "Monitor CPU utilization", Config: `{"threshold_percent": 80}`},
	// 		{Name: "Check memory usage", Description: "Monitor memory utilization", Config: `{"threshold_percent": 90}`},
	// 		{Name: "Check disk usage", Description: "Monitor disk utilization", Config: `{"threshold_percent": 85}`},
	// 		{Name: "Check service status", Description: "Check running services status", Config: `{"critical_services": ["nginx", "docker"]}`},
	// 		{Name: "Check log files", Description: "Analyze system log files", Config: `{"error_patterns": ["error", "critical", "failed"]}`},
	// 	},
	// }

	// // Create pipeline templates with their specific steps
	// pipelines := []entities.PipelineTemplate{
	// 	{
	// 		Name:        "System Update",
	// 		Description: "Update system packages and dependencies",
	// 	},
	// 	{
	// 		Name:        "Database Backup",
	// 		Description: "Perform a full backup of all databases",
	// 	},
	// 	{
	// 		Name:        "Network Security Scan",
	// 		Description: "Perform a comprehensive network security scan",
	// 	},
	// 	{
	// 		Name:        "Docker Container Cleanup",
	// 		Description: "Clean up unused Docker containers, images, and volumes",
	// 	},
	// 	{
	// 		Name:        "SSL Certificate Renewal",
	// 		Description: "Check and renew SSL certificates",
	// 	},
	// 	{
	// 		Name:        "System Health Check",
	// 		Description: "Perform a comprehensive system health check",
	// 	},
	// }

	// for _, pipeline := range pipelines {
	// 	if err := db.Get().Create(&pipeline).Error; err != nil {
	// 		logrus.Errorf("❌ Failed to create pipeline template %s: %v", pipeline.Name, err)
	// 		continue
	// 	}

	// 	steps := pipelineSteps[pipeline.Name]
	// 	for _, step := range steps {
	// 		if err := db.Get().Create(&step).Error; err != nil {
	// 			logrus.Errorf("❌ Failed to create step %s for pipeline %s: %v", step.Name, pipeline.Name, err)
	// 			continue
	// 		}
	// 	}

	// 	if err := db.Get().Model(&pipeline).Association("Steps").Append(steps); err != nil {
	// 		logrus.Errorf("❌ Failed to associate steps with pipeline %s: %v", pipeline.Name, err)
	// 		continue
	// 	}

	// 	logrus.Infof("✅ Created pipeline template: %s with %d steps", pipeline.Name, len(steps))
	// }
}

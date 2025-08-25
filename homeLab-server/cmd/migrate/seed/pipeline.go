package seed

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type ConfigFile struct {
	Providers []ProviderConfig `yaml:"providers"`
	Pipelines []PipelineConfig `yaml:"pipelines"`
}

type ProviderConfig struct {
	Name  string       `yaml:"name"`
	Steps []StepConfig `yaml:"steps"`
}

type StepConfig struct {
	Type        string                 `yaml:"type"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Config      map[string]interface{} `yaml:"config"`
}

type PipelineConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Steps       []string `yaml:"steps"`
}

// loadConfigFile reads and parses the YAML configuration file
func loadConfigFile(filePath string) (*ConfigFile, error) {
	yamlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()

	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	var config ConfigFile
	if err := yaml.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// findOrCreateStep handles finding or creating a step template
func findOrCreateStep(db *gorm.DB, repo repositories.PipelineRepository, stepConfig StepConfig, shouldUpdate bool) (*entities.StepTemplate, error) {
	// Direct database query for better performance
	var existingStep entities.StepTemplate
	result := db.Where("name = ?", stepConfig.Name).First(&existingStep)

	configJSON, err := json.Marshal(stepConfig.Config)
	if err != nil {
		return nil, err
	}

	if result.Error == nil {
		if shouldUpdate {
			existingStep.Type = stepConfig.Type
			existingStep.Description = stepConfig.Description
			existingStep.Config = string(configJSON)
			if err := repo.UpdateStepTemplate(&existingStep); err != nil {
				return nil, err
			}
			logrus.Infof("✅ Updated step template: %s", existingStep.Name)
		} else {
			logrus.Infof("⏩ Step %s already exists, skipping", stepConfig.Name)
		}
		return &existingStep, nil
	}

	step := &entities.StepTemplate{
		Type:        stepConfig.Type,
		Name:        stepConfig.Name,
		Description: stepConfig.Description,
		Config:      string(configJSON),
	}

	if err := repo.CreateStepTemplate(step); err != nil {
		return nil, err
	}
	logrus.Infof("✅ Created step template: %s", step.Name)
	return step, nil
}

// findOrCreatePipeline handles finding or creating a pipeline template
func findOrCreatePipeline(db *gorm.DB, repo repositories.PipelineRepository, pipelineConfig PipelineConfig, shouldUpdate bool) (*entities.PipelineTemplate, error) {
	// Direct database query for better performance
	var existingPipeline entities.PipelineTemplate
	result := db.Where("name = ?", pipelineConfig.Name).First(&existingPipeline)

	if result.Error == nil {
		if shouldUpdate {
			existingPipeline.Description = pipelineConfig.Description
			if err := repo.UpdatePipelineTemplate(&existingPipeline); err != nil {
				return nil, err
			}
			logrus.Infof("✅ Updated pipeline template: %s", existingPipeline.Name)
		} else {
			logrus.Infof("⏩ Pipeline %s already exists, skipping", pipelineConfig.Name)
		}
		return &existingPipeline, nil
	}

	pipeline := &entities.PipelineTemplate{
		Name:        pipelineConfig.Name,
		Description: pipelineConfig.Description,
	}

	if err := repo.CreatePipelineTemplate(pipeline); err != nil {
		return nil, err
	}
	logrus.Infof("✅ Created pipeline template: %s", pipeline.Name)
	return pipeline, nil
}

// associateStepsWithPipeline handles the association between pipeline and steps
func associateStepsWithPipeline(db *gorm.DB, repo repositories.PipelineRepository, pipeline *entities.PipelineTemplate, stepNames []string) error {
	// Direct database query for better performance
	var steps []entities.StepTemplate
	if err := db.Where("name IN ?", stepNames).Find(&steps).Error; err != nil {
		return err
	}

	// Collect step IDs
	stepIDs := make([]uint, len(steps))
	for i, step := range steps {
		stepIDs[i] = step.ID
	}

	// Associate steps with pipeline using repository
	if err := repo.AssociateStepsWithPipeline(pipeline.ID, stepIDs); err != nil {
		return err
	}

	logrus.Infof("✅ Associated %d steps with pipeline %s", len(stepIDs), pipeline.Name)
	return nil
}

func LoadDataFromYaml(db database.Database, shouldUpdate bool, config *ConfigFile) {
	// Initialize repository and get direct DB instance
	repo := repositories.NewPipelineRepository(db, nil, nil)
	gormDB := db.Get()

	// Process steps from all providers
	for _, provider := range config.Providers {
		logrus.Infof("📦 Processing provider: %s", provider.Name)
		for _, stepConfig := range provider.Steps {
			if _, err := findOrCreateStep(gormDB, repo, stepConfig, shouldUpdate); err != nil {
				logrus.Errorf("❌ Failed to process step %s: %v", stepConfig.Name, err)
			}
		}
	}

	// Process pipelines
	for _, pipelineConfig := range config.Pipelines {
		pipeline, err := findOrCreatePipeline(gormDB, repo, pipelineConfig, shouldUpdate)
		if err != nil {
			logrus.Errorf("❌ Failed to process pipeline %s: %v", pipelineConfig.Name, err)
			continue
		}

		if err := associateStepsWithPipeline(gormDB, repo, pipeline, pipelineConfig.Steps); err != nil {
			logrus.Errorf("❌ Failed to associate steps with pipeline %s: %v", pipeline.Name, err)
		}
	}
}

func CreatePipelineForSingleStep(db database.Database, config *ConfigFile) {
	repo := repositories.NewPipelineRepository(db, nil, nil)
	gormDB := db.Get()

	// Process steps from all providers
	for _, provider := range config.Providers {
		for _, stepConfig := range provider.Steps {
			pipelineConfig := PipelineConfig{
				Name:        fmt.Sprintf("[Single] %s", stepConfig.Name),
				Description: fmt.Sprintf("Single step pipeline for %s", stepConfig.Name),
				Steps:       []string{stepConfig.Name},
			}

			pipeline, err := findOrCreatePipeline(gormDB, repo, pipelineConfig, true)
			if err != nil {
				logrus.Errorf("❌ Failed to create single step pipeline for %s: %v", stepConfig.Name, err)
				continue
			}

			if err := associateStepsWithPipeline(gormDB, repo, pipeline, pipelineConfig.Steps); err != nil {
				logrus.Errorf("❌ Failed to associate step with single step pipeline %s: %v", pipeline.Name, err)
			}
		}
	}
}

func SeedPipelines(db database.Database) {
	config, err := loadConfigFile("pipeline.config.yaml")
	if err != nil {
		logrus.Errorf("❌ Failed to load pipeline.config.yaml: %v", err)
		return
	}
	LoadDataFromYaml(db, true, config)
	CreatePipelineForSingleStep(db, config)
}

package entities

import (
	"gorm.io/gorm"
)

type StatusEntity struct {
	gorm.Model
	Version string
	Health  string
}

type HealthEntity struct {
	gorm.Model
	Name    string
	IsOk    bool
	IsLive  bool
	Message string
}

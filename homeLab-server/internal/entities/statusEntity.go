package entities

import (
	"gorm.io/gorm"
)

type StatusEntity struct {
	gorm.Model
	Version string
	Health  string
}

package entities

import (
	"encoding/json"
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

type ServiceEntity struct {
	gorm.Model
	Name        string `json:"name"`
	Tags        string `json:"tags"`
	Description string `json:"description"`
	Namespace   string `json:"namespace"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	IP          string `json:"ip"`
	Port        int64  `json:"port"`
	LogoURL     string `json:"logo_url"`
}

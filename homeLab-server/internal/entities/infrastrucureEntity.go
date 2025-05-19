package entities

type Infrastructure struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	IPAddress string `gorm:"not null"`
	Status    string `gorm:"not null"`
}

type InfrastructureStatus string

const (
	StatusActive   InfrastructureStatus = "active"
	StatusInactive InfrastructureStatus = "inactive"
)

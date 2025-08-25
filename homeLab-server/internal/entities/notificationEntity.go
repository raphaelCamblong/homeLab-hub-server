package entities

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func UnmarshalNotification(data []byte) (Notification, error) {
	var r Notification
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Notification) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Notification struct {
	gorm.Model

	// Required core fields
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`     // info, warning, error, success
	Severity  string    `json:"severity" gorm:"not null"` // low, medium, high
	Source    string    `json:"source" gorm:"not null"`   // service name that generated the notification
	CreatedAt time.Time `json:"created_at" gorm:"not null"`

	// Optional
	Tags     string     `json:"tags"`
	Metadata string     `json:"metadata"` // Additional JSON data
	IsRead   bool       `json:"is_read" gorm:"default:false"`
	ReadAt   *time.Time `json:"read_at"`
}

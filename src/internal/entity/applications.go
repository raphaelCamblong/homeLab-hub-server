package entity

import "time"

type Application struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags"`
	LogoPath    string       `json:"logo_path,omitempty"`
	State       ServiceState `json:"state,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
}

type ServiceState string

const (
	Running ServiceState = "running"
	Stopped ServiceState = "stopped"
)

func NewApplication(name string, description string, tags []string, logoPath string, state ServiceState) *Application {
	return &Application{Name: name, Description: description, Tags: tags, LogoPath: logoPath, State: state}
}

package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type SystemUseCase interface {
	GetSystemInfo() (*SystemInfo, error)
	ListProcesses() ([]Process, error)
	ListActiveServices() ([]Service, error)
	FetchSystemMetrics() (*Metrics, error)
	CheckHealth() (*HealthStatus, error)
	RestartSystem() error
	ShutdownSystem() error
}

type systemUseCase struct {
	systemRepository repositories.SystemRepository
}

func NewSystemUseCase(systemRepository repositories.SystemRepository) SystemUseCase {
	return &systemUseCase{systemRepository: systemRepository}
}

func (u *systemUseCase) GetSystemInfo() (*SystemInfo, error) {
	return u.systemRepository.GetSystemInfo()
}

func (u *systemUseCase) ListProcesses() ([]Process, error) {
	return u.systemRepository.ListProcesses()
}

func (u *systemUseCase) ListActiveServices() ([]Service, error) {
	return u.systemRepository.ListActiveServices()
}

func (u *systemUseCase) FetchSystemMetrics() (*Metrics, error) {
	return u.systemRepository.FetchSystemMetrics()
}

func (u *systemUseCase) CheckHealth() (*HealthStatus, error) {
	return u.systemRepository.CheckHealth()
}

func (u *systemUseCase) RestartSystem() error {
	return u.systemRepository.RestartSystem()
}

func (u *systemUseCase) ShutdownSystem() error {
	return u.systemRepository.ShutdownSystem()
} 
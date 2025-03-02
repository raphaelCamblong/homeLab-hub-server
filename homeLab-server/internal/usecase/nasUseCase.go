package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type NasUseCase interface {
	ListStoragePools() ([]StoragePool, error)
	GetStoragePoolDetails(id string) (*StoragePoolDetails, error)
	GetZFSMetrics() (*ZFSMetrics, error)
	BrowseFilesystem() ([]File, error)
	UploadFile(file *File) error
	DownloadFile(fileID string) (*File, error)
	TriggerBackupJob() error
}

type nasUseCase struct {
	nasRepository repositories.NASRepository
}

func NewNasUseCase(nasRepository repositories.NASRepository) NasUseCase {
	return &nasUseCase{nasRepository: nasRepository}
}

func (u *nasUseCase) ListStoragePools() ([]StoragePool, error) {
	return u.nasRepository.ListStoragePools()
}

func (u *nasUseCase) GetStoragePoolDetails(id string) (*StoragePoolDetails, error) {
	return u.nasRepository.GetStoragePoolDetails(id)
}

func (u *nasUseCase) GetZFSMetrics() (*ZFSMetrics, error) {
	return u.nasRepository.GetZFSMetrics()
}

func (u *nasUseCase) BrowseFilesystem() ([]File, error) {
	return u.nasRepository.BrowseFilesystem()
}

func (u *nasUseCase) UploadFile(file *File) error {
	return u.nasRepository.UploadFile(file)
}

func (u *nasUseCase) DownloadFile(fileID string) (*File, error) {
	return u.nasRepository.DownloadFile(fileID)
}

func (u *nasUseCase) TriggerBackupJob() error {
	return u.nasRepository.TriggerBackupJob()
} 
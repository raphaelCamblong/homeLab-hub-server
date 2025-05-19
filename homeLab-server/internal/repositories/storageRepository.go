package repositories

type StorageRepository interface{}

type storageRepository struct{}

func NewStorageRepository() StorageRepository {
	return &storageRepository{}
}

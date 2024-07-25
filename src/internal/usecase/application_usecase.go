package usecase

import (
	"example.com/service-api/src/internal/entity"
)

type ApplicationRepository interface {
	Save(application *entity.Application) error
	FindById(id string) (*entity.Application, error)
}

type ApplicationUseCase struct {
	Repo ApplicationRepository
}

func (uc *ApplicationUseCase) CreateApplication(name string, description string, tags []string, logoPath string, state ServiceState) error {
	application := entity.NewApplication(name, description, tags, logoPath, state)
	return uc.Repo.Save(application)
}

func (uc *ApplicationUseCase) FindById(id string) (*entity.Application, error) {
	application, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return application, nil
}

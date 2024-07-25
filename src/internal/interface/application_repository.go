package _interface

import "example.com/service-api/src/internal/entity"

type ApplicationRepository interface {
	Save(application *entity.Application) error
	FindById(id string) (*entity.Application, error)
}

package entity

import "context"

type ServiceRepository interface {
	CreateService(ctx context.Context, service Service) (Service, error)
}

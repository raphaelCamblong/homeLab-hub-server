package repositories

import (
	"encoding/json"

	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type IotRepository interface {
	GetDevices() ([]entities.IotPlugDevice, error)
}

type iotRepository struct {
	lambdaClient client.LambdaFunction
}

func NewIotRepository(lambdaClient client.LambdaFunction) IotRepository {
	return &iotRepository{
		lambdaClient: lambdaClient,
	}
}

func (r *iotRepository) GetDevices() ([]entities.IotPlugDevice, error) {
	cfg := client.LambdaFnConfig{
		FunctionName: "network-controller",
		Args: map[string]interface{}{
			"action": "get-devices",
		},
	}

	bytes, err := r.lambdaClient.Invoke(cfg)
	if err != nil {
		return nil, err
	}

	var devices []entities.IotPlugDevice
	if err := json.Unmarshal(bytes, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

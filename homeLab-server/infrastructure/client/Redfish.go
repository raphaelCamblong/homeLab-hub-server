package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"io"
	"net/http"
	"homelab.com/homelab-server/homeLab-server/internal/entities/ilo"
)

type Redfish interface {
	CreateSession(credentials *Credentials) (*RequestOption, error)
	GetThermalData(*RequestOption) (*ilo.ThermalEntity, error)
	GetPowerFastData(*RequestOption) (*ilo.PowerEntity, error)
	GetPowerData(*RequestOption) (*ilo.PowerEntity, error)
	GetHealthCheck(*RequestOption) (*ilo.HealthEntity, error)
}

type Credentials struct {
	Username string `json:"UserName"`
	Password string `json:"Password"`
}

type redfish struct {
	cfg config.NetworkConnection
}

func NewRedfishInfra(c config.NetworkConnection) Redfish {
	return &redfish{c}
}

func (r *redfish) CreateSession(cred *Credentials) (*RequestOption, error) {
	sessionURL := fmt.Sprintf("%s/Sessions", r.cfg.Host)
	credJson, err := json.Marshal(cred)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	reqBody, err := io.ReadAll(bytes.NewReader(credJson))
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON data: %w", err)
	}

	resp, err := http.Post(sessionURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create Redfish session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create Redfish session: status code %d", resp.StatusCode)
	}

	token := resp.Header.Get("X-Auth-Token")
	if token == "" {
		return nil, fmt.Errorf("failed to decode Redfish session token header response empty")
	}
	return &RequestOption{token}, nil
}

func (r *redfish) getData(path string, requestCtx *RequestOption) (*[]byte, error) {
	thermalURL := fmt.Sprintf("%s/%s", r.cfg.Host, path)
	req, err := http.NewRequest(http.MethodGet, thermalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Auth-Token", requestCtx.AuthToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve thermal data: status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &bodyBytes, nil
}

func (r *redfish) GetThermalData(requestCtx *RequestOption) (*ilo.ThermalEntity, error) {
	res, err := r.getData("Chassis/1/Thermal", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}

	var thermalData ilo.ThermalEntity
	if err := json.Unmarshal(*res, &thermalData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal thermal data: %w", err)
	}
	return &thermalData, nil
}

func (r *redfish) GetPowerData(requestCtx *RequestOption) (*ilo.PowerEntity, error) {
	res, err := r.getData("Chassis/1/Power/PowerMeter", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power data: %w", err)
	}

	var powerData ilo.PowerEntity
	if err := json.Unmarshal(*res, &powerData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power data: %w", err)
	}
	return &powerData, nil
}

func (r *redfish) GetPowerFastData(requestCtx *RequestOption) (*ilo.PowerEntity, error) {
	res, err := r.getData("Chassis/1/Power/FastPowerMeter", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power fast data: %w", err)
	}

	var powerFastData ilo.PowerEntity
	if err := json.Unmarshal(*res, &powerFastData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power fast data: %w", err)
	}
	return &powerFastData, nil
}

func (r *redfish) GetHealthCheck(requestCtx *RequestOption) (*ilo.HealthEntity, error) {
	res, err := r.getData("Health", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve health data: %w", err)
	}

	var healthData ilo.HealthEntity
	if err := json.Unmarshal(*res, &healthData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal health data: %w", err)
	}
	return &healthData, nil
}

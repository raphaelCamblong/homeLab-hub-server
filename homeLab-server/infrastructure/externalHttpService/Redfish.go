package externalHttpService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Redfish interface {
	CreateSession(credentials *Credentials) (*RequestOption, error)
	GetThermalData(request *RequestOption) (*[]byte, error)
	GetPowerFastData(request *RequestOption) (*[]byte, error)
	GetPowerData(request *RequestOption) (*[]byte, error)
}

type Credentials struct {
	Username string `json:"UserName"`
	Password string `json:"Password"`
}

type redfish struct {
	BaseUrl string
}

func NewRedfishInfra(baseUrl string) Redfish {
	return &redfish{baseUrl}
}

func (r *redfish) CreateSession(cred *Credentials) (*RequestOption, error) {
	sessionURL := fmt.Sprintf("%s/Sessions", r.BaseUrl)
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
	thermalURL := fmt.Sprintf("%s/%s", r.BaseUrl, path)
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

func (r *redfish) GetThermalData(requestCtx *RequestOption) (*[]byte, error) {
	res, err := r.getData("Chassis/1/Thermal", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return res, nil
}

func (r *redfish) GetPowerData(requestCtx *RequestOption) (*[]byte, error) {
	res, err := r.getData("Chassis/1/Power/PowerMeter", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power data: %w", err)
	}
	return res, nil
}

func (r *redfish) GetPowerFastData(requestCtx *RequestOption) (*[]byte, error) {
	res, err := r.getData("Chassis/1/Power/FastPowerMeter", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve power fast data: %w", err)
	}
	return res, nil
}

package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"io"
	"net/http"
	"time"
)

type RedfishRepository interface {
	CreateSession() error
	IsSessionOpen() bool
	getTokenFromServer() (string, error)
	getCachedToken() (string, error)
	saveToken() error
	GetToken() (string, error)
	GetThermalData() (*entities.ThermalEntity, error)
}

type redfishRepository struct {
	BaseUrl string
	token   string
	Cache   cache.Database
}

type Credentials struct {
	Username string `json:"UserName"`
	Password string `json:"Password"`
}

func NewRedfishRepository(cache cache.Database) RedfishRepository {
	return &redfishRepository{
		BaseUrl: "",
		token:   "",
		Cache:   cache,
	}
}

func (r *redfishRepository) CreateSession() error {
	cfg := config.GetConfig()
	//TODO: maybe move this to infra
	r.BaseUrl = fmt.Sprintf("https://%s/rest/v1", cfg.ExternalServicesCredential.IloIp)

	if !r.IsSessionOpen() {
		if _, err := r.getTokenFromServer(); err != nil {
			return fmt.Errorf("unable to get token from server: %w", err)
		}
		//if err := r.saveToken(); err != nil {
		//	return err
		//}
	}
	return nil
}

func (r *redfishRepository) IsSessionOpen() bool {
	return r.token != ""
}

func (r *redfishRepository) getTokenFromServer() (string, error) {
	cfg := config.GetConfig()
	cred := Credentials{
		Username: cfg.ExternalServicesCredential.IloUsername,
		Password: cfg.ExternalServicesCredential.IloPassword,
	}

	sessionURL := fmt.Sprintf("%s/Sessions", r.BaseUrl)
	credJson, err := json.Marshal(cred)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %w", err)
	}

	reqBody, err := io.ReadAll(bytes.NewReader(credJson))
	if err != nil {
		return "", fmt.Errorf("failed to read JSON data: %w", err)
	}

	resp, err := http.Post(sessionURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create Redfish session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create Redfish session: status code %d", resp.StatusCode)
	}

	token := resp.Header.Get("X-Auth-Token")

	if token == "" {
		return "", fmt.Errorf("failed to decode Redfish session token header response empty")
	}
	r.token = token
	fmt.Printf("Get token from server: %s", r.token)
	return r.token, nil
}

func (r *redfishRepository) getCachedToken() (string, error) {
	ctx := context.Background()

	val, err := r.Cache.GetClient().Get(ctx, "redfishToken").Result()
	if err != nil {
		return "", fmt.Errorf("can't get Redfish cached Token")
	}
	r.token = val
	return val, nil
}

func (r *redfishRepository) saveToken() error {
	ctx := context.Background()

	err := r.Cache.GetClient().Set(ctx, "redfishToken", r.token, time.Minute*30).Err()
	if err != nil {
		return fmt.Errorf("can't save Redfish Token")
	}
	return nil
}

func (r *redfishRepository) GetToken() (string, error) {
	if r.IsSessionOpen() {
		return r.token, nil
	}
	return r.getCachedToken()
}

func (r *redfishRepository) GetThermalData() (*entities.ThermalEntity, error) {
	token, _ := r.GetToken()
	thermalURL := fmt.Sprintf("%s/Chassis/1/Thermal", r.BaseUrl)
	req, err := http.NewRequest(http.MethodGet, thermalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Auth-Token", token)
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

	var thermalData entities.ThermalEntity
	thermalData, err = entities.UnmarshalThermalEntity(bodyBytes)

	return &thermalData, err
}

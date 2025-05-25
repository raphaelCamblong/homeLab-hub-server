package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"homelab.com/homelab-server/homeLab-server/init/config"
)

type LambdaFunction interface {
	Invoke(fnCfg LambdaFnConfig) ([]byte, error)
}

type LambdaFnConfig struct {
	FunctionName string                 `json:"functionName"`
	BaseURL      string                 `json:"baseUrl,omitempty"`
	Token        string                 `json:"token,omitempty"`
	Args         map[string]interface{} `json:"args,omitempty"`
	Headers      map[string]string      `json:"headers,omitempty"`
}

func NewLambdaFnConfigUnmarshal(cfg string) (*LambdaFnConfig, error) {
	var lambdaFnConfig map[string]interface{}
	if err := json.Unmarshal([]byte(cfg), &lambdaFnConfig); err != nil {
		return nil, err
	}
	return NewLambdaFnConfig(lambdaFnConfig)
}

func NewLambdaFnConfig(cfg map[string]interface{}) (*LambdaFnConfig, error) {
	config := &LambdaFnConfig{}

	if err := validateAndSetFunctionName(config, cfg); err != nil {
		return nil, err
	}

	setOptionalConfigs(config, cfg)
	return config, nil
}

func validateAndSetFunctionName(config *LambdaFnConfig, cfg map[string]interface{}) error {
	if fnName, ok := cfg["fn_name"].(string); ok {
		config.FunctionName = fnName
		return nil
	}
	return fmt.Errorf("invalid or missing function name")
}

func setOptionalConfigs(config *LambdaFnConfig, cfg map[string]interface{}) {
	if baseURL, ok := cfg["base_url"].(string); ok {
		config.BaseURL = baseURL
	}

	if token, ok := cfg["token"].(string); ok {
		config.Token = token
	}

	if args, ok := cfg["args"].(map[string]interface{}); ok {
		config.Args = args
	}

	if headers, ok := cfg["headers"].(map[string]string); ok {
		config.Headers = headers
	}
}

type lambdaFunction struct {
	cfg    config.LambdaConfig
	client HttpClient
}

func NewLambdaFunction(cfg config.LambdaConfig, client HttpClient) LambdaFunction {
	return &lambdaFunction{
		cfg:    cfg,
		client: client,
	}
}

func (l *lambdaFunction) Invoke(fnCfg LambdaFnConfig) ([]byte, error) {
	url, err := l.buildURL(fnCfg)
	if err != nil {
		return nil, err
	}

	payload, err := l.preparePayload(fnCfg)
	if err != nil {
		return nil, err
	}

	headers := l.prepareHeaders(fnCfg)

	resp, err := l.makeRequest(url, payload, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return l.handleResponse(resp)
}

func (l *lambdaFunction) buildURL(fnCfg LambdaFnConfig) (string, error) {
	baseURL := fnCfg.BaseURL
	if baseURL == "" {
		baseURL = l.cfg.BaseURL
	}
	return fmt.Sprintf("%s/%s", baseURL, fnCfg.FunctionName), nil
}

func (l *lambdaFunction) preparePayload(fnCfg LambdaFnConfig) ([]byte, error) {
	if fnCfg.Args == nil {
		return []byte("{}"), nil
	}

	payload, err := json.Marshal(fnCfg.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal args: %w", err)
	}
	return payload, nil
}

func (l *lambdaFunction) prepareHeaders(fnCfg LambdaFnConfig) map[string]string {
	headers := make(map[string]string)

	headers["Content-Type"] = "application/json"

	if fnCfg.Token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", fnCfg.Token)
	}

	for key, value := range fnCfg.Headers {
		headers[key] = value
	}

	return headers
}

func (l *lambdaFunction) makeRequest(url string, payload []byte, headers map[string]string) (*http.Response, error) {
	resp, err := l.client.Get(context.Background(), url, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke function: %w", err)
	}
	return resp, nil
}

func (l *lambdaFunction) handleResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("function returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

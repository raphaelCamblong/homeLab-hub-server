package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"homelab.com/homelab-server/homeLab-server/init/config"
)

type LambdaFunction interface {
	Invoke(functionName string, payload []byte) ([]byte, error)
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

func (l *lambdaFunction) Invoke(functionName string, payload []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/function/%s", l.cfg.BaseURL, functionName)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Post(context.Background(), url, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke function: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("function returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

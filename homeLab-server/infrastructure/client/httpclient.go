package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(timeout time.Duration) HttpClient {
	return &httpClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *httpClient) Get(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	bodyBytes, err := c.marshalBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addHeaders(req, headers)
	return c.Do(req)
}

func (c *httpClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	bodyBytes, err := c.marshalBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addHeaders(req, headers)
	return c.Do(req)
}

func (c *httpClient) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	bodyBytes, err := c.marshalBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addHeaders(req, headers)
	return c.Do(req)
}

func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addHeaders(req, headers)
	return c.Do(req)
}

func (c *httpClient) addHeaders(req *http.Request, headers map[string]string) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = "application/json"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func (c *httpClient) marshalBody(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch v := body.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case io.Reader:
		return io.ReadAll(v)
	default:
		return json.Marshal(v)
	}
}

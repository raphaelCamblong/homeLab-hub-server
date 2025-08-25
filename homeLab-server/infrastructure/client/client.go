package client

import (
	"log"
	"time"
)

type ExternalHttpService interface {
	GetK8sClient() *K8sClient
	GetHttpClient() HttpClient
}

type client struct {
	httpClient HttpClient
	k8sClient  *K8sClient
}

func NewExternalHttpService() ExternalHttpService {
	k8sClient, err := NewK8sClient(&ClientOptions{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}
	return &client{
		httpClient: NewHttpClient(120 * time.Second),
		k8sClient:  k8sClient,
	}
}

func (c *client) GetK8sClient() *K8sClient {
	return c.k8sClient
}

func (c *client) GetHttpClient() HttpClient {
	return c.httpClient
}

package client

import "time"

type ExternalHttpService interface {
}

type client struct {
	httpClient HttpClient
}

func NewExternalHttpService() ExternalHttpService {
	return &client{
		httpClient: NewHttpClient(120 * time.Second),
	}
}

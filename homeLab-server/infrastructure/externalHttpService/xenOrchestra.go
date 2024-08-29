package externalHttpService

import (
	"fmt"
	"io"
	"net/http"
)

type XenOrchestra interface {
	GetAllVm(*RequestOption) (*[]byte, error)
	GetVm(string, *RequestOption) (*[]byte, error)
	GetAllHost(*RequestOption) (*[]byte, error)
	GetHost(string, *RequestOption) (*[]byte, error)
	getData(string, *RequestOption) (*[]byte, error)
}

type (
	xenOrchestra struct {
		BaseUrl string
	}
)

func NewXenOrchestraInfra(baseUrl string) XenOrchestra {
	return &xenOrchestra{baseUrl}
}

func (x *xenOrchestra) GetAllVm(requestCtx *RequestOption) (*[]byte, error) {
	res, err := x.getData("/vms", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return res, nil
}

func (x *xenOrchestra) GetVm(id string, requestCtx *RequestOption) (*[]byte, error) {
	slug := fmt.Sprintf("/vms/%s", id)
	res, err := x.getData(slug, requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return res, nil
}

func (x *xenOrchestra) GetAllHost(requestCtx *RequestOption) (*[]byte, error) {
	res, err := x.getData("/hosts", requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return res, nil
}

func (x *xenOrchestra) GetHost(id string, requestCtx *RequestOption) (*[]byte, error) {
	slug := fmt.Sprintf("/hosts/%s", id)
	res, err := x.getData(slug, requestCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w", err)
	}
	return res, nil
}

func (x *xenOrchestra) getData(path string, requestCtx *RequestOption) (*[]byte, error) {
	var separator string = ""
	if path[0] != '/' {
		separator = "/"
	}
	url := fmt.Sprintf("%s%s%s", x.BaseUrl, separator, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "authenticationToken",
		Value: requestCtx.AuthToken,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thermal data: %w %s", err, url)
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

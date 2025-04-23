package client

type ExternalHttpService interface {
	GetRedfish() Redfish
	GetXenOrchestra() XenOrchestra
	GetPrometheus() Prometheus
}

type client struct {
	redfish   Redfish
	xo        XenOrchestra
	prometheus Prometheus
}

func NewExternalHttpService(redfish Redfish, xen XenOrchestra, prometheus Prometheus) ExternalHttpService {
	return &client{redfish, xen, prometheus}
}

func (e *client) GetRedfish() Redfish {
	return e.redfish
}

func (e *client) GetXenOrchestra() XenOrchestra {
	return e.xo
}

func (e *client) GetPrometheus() Prometheus {
	return e.prometheus
}

type (
	RequestOption struct {
		AuthToken string `json:"AuthToken"`
	}
	AuthToken struct {
		AuthToken string `json:"authenticationToken"`
	}
)

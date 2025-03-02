package externalHttpService

type ExternalHttpService interface {
	GetRedfish() Redfish
	GetXenOrchestra() XenOrchestra
	GetPrometheus() Prometheus
}

type externalHttpService struct {
	redfish   Redfish
	xo        XenOrchestra
	prometheus Prometheus
}

func NewExternalHttpService(redfish Redfish, xen XenOrchestra, prometheus Prometheus) ExternalHttpService {
	return &externalHttpService{redfish, xen, prometheus}
}

func (e *externalHttpService) GetRedfish() Redfish {
	return e.redfish
}

func (e *externalHttpService) GetXenOrchestra() XenOrchestra {
	return e.xo
}

func (e *externalHttpService) GetPrometheus() Prometheus {
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

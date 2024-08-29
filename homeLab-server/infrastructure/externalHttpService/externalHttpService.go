package externalHttpService

type ExternalHttpService interface {
	GetRedfish() Redfish
	GetXenOrchestra() XenOrchestra
}

type externalHttpService struct {
	redfish Redfish
	xo      XenOrchestra
}

func NewExternalHttpService(redfish Redfish, xen XenOrchestra) ExternalHttpService {
	return &externalHttpService{redfish, xen}
}

func (e *externalHttpService) GetRedfish() Redfish {
	return e.redfish
}

func (e *externalHttpService) GetXenOrchestra() XenOrchestra {
	return e.xo
}

type (
	RequestOption struct {
		AuthToken string `json:"AuthToken"`
	}
	AuthToken struct {
		AuthToken string `json:"authenticationToken"`
	}
)

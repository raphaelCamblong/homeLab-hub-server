package externalHttpService

type ExternalHttpService interface {
	GetRedfish() Redfish
	GetXenOrchestra() any
}

type externalHttpService struct {
	redfish Redfish
	xo      any
}

func NewExternalHttpService(redfish Redfish, xen any) ExternalHttpService {
	return &externalHttpService{redfish, xen}
}

func (e *externalHttpService) GetRedfish() Redfish {
	return e.redfish
}

func (e *externalHttpService) GetXenOrchestra() any {
	return e.xo
}

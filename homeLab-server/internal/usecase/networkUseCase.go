package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type NetworkUseCase interface {
	ListDevices() ([]Device, error)
	ListNetworkInterfaces() ([]NetworkInterface, error)
	GetInterfaceDetails(id string) (*NetworkInterfaceDetails, error)
	ViewFirewallRules() ([]FirewallRule, error)
	ModifyFirewallRules(rules []FirewallRule) error
	ViewDNSSettings() (*DNSSettings, error)
	RunNetworkSpeedTest() (*SpeedTestResult, error)
}

type networkUseCase struct {
	networkRepository repositories.NetworkRepository
}

func NewNetworkUseCase(networkRepository repositories.NetworkRepository) NetworkUseCase {
	return &networkUseCase{networkRepository: networkRepository}
}

func (u *networkUseCase) ListDevices() ([]Device, error) {
	return u.networkRepository.ListDevices()
}

func (u *networkUseCase) ListNetworkInterfaces() ([]NetworkInterface, error) {
	return u.networkRepository.ListNetworkInterfaces()
}

func (u *networkUseCase) GetInterfaceDetails(id string) (*NetworkInterfaceDetails, error) {
	return u.networkRepository.GetInterfaceDetails(id)
}

func (u *networkUseCase) ViewFirewallRules() ([]FirewallRule, error) {
	return u.networkRepository.ViewFirewallRules()
}

func (u *networkUseCase) ModifyFirewallRules(rules []FirewallRule) error {
	return u.networkRepository.ModifyFirewallRules(rules)
}

func (u *networkUseCase) ViewDNSSettings() (*DNSSettings, error) {
	return u.networkRepository.ViewDNSSettings()
}

func (u *networkUseCase) RunNetworkSpeedTest() (*SpeedTestResult, error) {
	return u.networkRepository.RunNetworkSpeedTest()
} 
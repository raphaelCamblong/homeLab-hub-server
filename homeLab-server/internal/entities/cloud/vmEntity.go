package entities

import "time"

import "encoding/json"

func UnmarshalVMEntity(data []byte) (VMEntity, error) {
	var r VMEntity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *VMEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VMEntity struct {
	Type                    string        `json:"type"`
	Addresses               Addresses     `json:"addresses"`
	AutoPoweron             bool          `json:"auto_poweron"`
	BIOSStrings             BIOSStrings   `json:"bios_strings"`
	BlockedOperations       Operations    `json:"blockedOperations"`
	Boot                    Boot          `json:"boot"`
	CPUs                    CPUs          `json:"CPUs"`
	Creation                Creation      `json:"creation"`
	CurrentOperations       Operations    `json:"current_operations"`
	ExpNestedHvm            bool          `json:"expNestedHvm"`
	Viridian                bool          `json:"viridian"`
	MainIPAddress           string        `json:"mainIpAddress"`
	HighAvailability        string        `json:"high_availability"`
	IsFirmwareSupported     bool          `json:"isFirmwareSupported"`
	Memory                  Memory        `json:"memory"`
	InstallTime             int64         `json:"installTime"`
	NameDescription         string        `json:"name_description"`
	NameLabel               string        `json:"name_label"`
	NeedsVtpm               bool          `json:"needsVtpm"`
	Other                   Other         `json:"other"`
	OSVersion               OSVersion     `json:"os_version"`
	PowerState              string        `json:"power_state"`
	HasVendorDevice         bool          `json:"hasVendorDevice"`
	Snapshots               []interface{} `json:"snapshots"`
	StartDelay              int64         `json:"startDelay"`
	StartTime               int64         `json:"startTime"`
	SecureBoot              bool          `json:"secureBoot"`
	Tags                    []interface{} `json:"tags"`
	VIFS                    []string      `json:"VIFs"`
	VTPMS                   []interface{} `json:"VTPMs"`
	VirtualizationMode      string        `json:"virtualizationMode"`
	XenTools                XenTools      `json:"xenTools"`
	ManagementAgentDetected bool          `json:"managementAgentDetected"`
	PVDriversDetected       bool          `json:"pvDriversDetected"`
	PVDriversVersion        string        `json:"pvDriversVersion"`
	PVDriversUpToDate       bool          `json:"pvDriversUpToDate"`
	Container               string        `json:"$container"`
	VBDs                    []string      `json:"$VBDs"`
	VMEntityVGPUs           []interface{} `json:"VGPUs"`
	VGPUs                   []interface{} `json:"$VGPUs"`
	XenStoreData            XenStoreData  `json:"xenStoreData"`
	VGA                     string        `json:"vga"`
	Videoram                string        `json:"videoram"`
	ID                      string        `json:"id"`
	UUID                    string        `json:"uuid"`
	Pool                    string        `json:"$pool"`
	PoolID                  string        `json:"$poolId"`
	XapiRef                 string        `json:"_xapiRef"`
	MessagesHref            string        `json:"messages_href"`
	VdisHref                string        `json:"vdis_href"`
}

type Addresses struct {
	The0Ipv40 string `json:"0/ipv4/0"`
	The0Ipv60 string `json:"0/ipv6/0"`
	The0Ipv61 string `json:"0/ipv6/1"`
}

type BIOSStrings struct {
	BIOSVendor                 *string `json:"bios-vendor"`
	BIOSVersion                *string `json:"bios-version"`
	SystemManufacturer         *string `json:"system-manufacturer"`
	SystemProductName          *string `json:"system-product-name"`
	SystemVersion              *string `json:"system-version"`
	SystemSerialNumber         *string `json:"system-serial-number"`
	BaseboardManufacturer      *string `json:"baseboard-manufacturer"`
	BaseboardProductName       *string `json:"baseboard-product-name"`
	BaseboardVersion           *string `json:"baseboard-version"`
	BaseboardSerialNumber      *string `json:"baseboard-serial-number"`
	BaseboardAssetTag          *string `json:"baseboard-asset-tag"`
	BaseboardLocationInChassis *string `json:"baseboard-location-in-chassis"`
	EnclosureAssetTag          *string `json:"enclosure-asset-tag"`
	OEM1                       string  `json:"oem-1"`
	OEM2                       string  `json:"oem-2"`
	OEM3                       string  `json:"oem-3"`
	OEM4                       string  `json:"oem-4"`
	OEM5                       string  `json:"oem-5"`
	HPRombios                  string  `json:"hp-rombios"`
}

type Operations struct {
}

type Boot struct {
	Firmware string `json:"firmware"`
	Order    string `json:"order"`
}

type CPUs struct {
	Max             *int64  `json:"max"`
	Number          *int64  `json:"number"`
	CPUCount        *string `json:"cpu_count"`
	SocketCount     *string `json:"socket_count"`
	Vendor          *string `json:"vendor"`
	Speed           *string `json:"speed"`
	Modelname       *string `json:"modelname"`
	Family          *string `json:"family"`
	Model           *string `json:"model"`
	Stepping        *string `json:"stepping"`
	Flags           *string `json:"flags"`
	FeaturesPV      *string `json:"features_pv"`
	FeaturesHvm     *string `json:"features_hvm"`
	FeaturesHvmHost *string `json:"features_hvm_host"`
	FeaturesPVHost  *string `json:"features_pv_host"`
}
type Creation struct {
	Date     time.Time `json:"date"`
	Template string    `json:"template"`
	User     string    `json:"user"`
}

type Memory struct {
	Dynamic *[]int64 `json:"dynamic"`
	Static  []int64  `json:"static"`
	Size    int64    `json:"size"`
}

type OSVersion struct {
	Name   string `json:"name"`
	Uname  string `json:"uname"`
	Distro string `json:"distro"`
	Major  string `json:"major"`
	Minor  string `json:"minor"`
}

type Other struct {
	Xo6897F0D5       string `json:"xo:6897f0d5"`
	AutoPoweron      string `json:"auto_poweron"`
	Xo0Df69429       string `json:"xo:0df69429"`
	Xo41B2E070       string `json:"xo:41b2e070"`
	BaseTemplateName string `json:"base_template_name"`
	ImportTask       string `json:"import_task"`
	MACSeed          string `json:"mac_seed"`
	InstallMethods   string `json:"install-methods"`
	LinuxTemplate    string `json:"linux_template"`
}

type XenStoreData struct {
	VMDataMMIOHoleSize string `json:"vm-data/mmio-hole-size"`
	VMData             string `json:"vm-data"`
}

type XenTools struct {
	Major   int64   `json:"major"`
	Minor   int64   `json:"minor"`
	Version float64 `json:"version"`
}

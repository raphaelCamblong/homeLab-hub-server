package entities

import "encoding/json"

func UnmarshalHostEntity(data []byte) (*HostEntity, error) {
	var r HostEntity
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *HostEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type HostEntity struct {
	CPUs               CPUs              `json:"CPUs"`
	Address            string            `json:"address"`
	BIOSStrings        BIOSStrings       `json:"bios_strings"`
	Build              string            `json:"build"`
	ChipsetInfo        ChipsetInfo       `json:"chipset_info"`
	Enabled            bool              `json:"enabled"`
	ControlDomain      string            `json:"controlDomain"`
	Cpus               Cpus              `json:"cpus"`
	CurrentOperations  CurrentOperations `json:"current_operations"`
	Hostname           string            `json:"hostname"`
	IscsiIqn           string            `json:"iscsiIqn"`
	ZstdSupported      bool              `json:"zstdSupported"`
	LicenseParams      map[string]string `json:"license_params"`
	LicenseServer      LicenseServer     `json:"license_server"`
	LicenseExpiry      interface{}       `json:"license_expiry"`
	Logging            CurrentOperations `json:"logging"`
	NameDescription    string            `json:"name_description"`
	NameLabel          string            `json:"name_label"`
	Memory             Memory            `json:"memory"`
	Multipathing       bool              `json:"multipathing"`
	OtherConfig        OtherConfig       `json:"otherConfig"`
	Patches            []interface{}     `json:"patches"`
	PowerOnMode        string            `json:"powerOnMode"`
	PowerState         string            `json:"power_state"`
	ResidentVms        []string          `json:"residentVms"`
	StartTime          int64             `json:"startTime"`
	SupplementalPacks  []interface{}     `json:"supplementalPacks"`
	AgentStartTime     int64             `json:"agentStartTime"`
	RebootRequired     bool              `json:"rebootRequired"`
	Tags               []interface{}     `json:"tags"`
	Version            string            `json:"version"`
	ProductBrand       string            `json:"productBrand"`
	HvmCapable         bool              `json:"hvmCapable"`
	Certificates       []interface{}     `json:"certificates"`
	HostEntityPIFS     []string          `json:"PIFs"`
	PIFS               []string          `json:"$PIFs"`
	HostEntityPCIs     []string          `json:"PCIs"`
	PCIs               []string          `json:"$PCIs"`
	HostEntityPGPUs    []string          `json:"PGPUs"`
	PGPUs              []string          `json:"$PGPUs"`
	PBDs               []string          `json:"$PBDs"`
	ID                 string            `json:"id"`
	Type               string            `json:"type"`
	UUID               string            `json:"uuid"`
	Pool               string            `json:"$pool"`
	PoolID             string            `json:"$poolId"`
	XapiRef            string            `json:"_xapiRef"`
	MessagesHref       string            `json:"messages_href"`
	AuditHref          string            `json:"audit_href"`
	LogsHref           string            `json:"logs_href"`
	MissingPatchesHref string            `json:"missing_patches_href"`
	SMTHref            string            `json:"smt_href"`
}

type ChipsetInfo struct {
	Iommu bool `json:"iommu"`
}

type Cpus struct {
	Cores   int64 `json:"cores"`
	Sockets int64 `json:"sockets"`
}

type CurrentOperations struct {
}

type LicenseServer struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type OtherConfig struct {
	AgentStartTime           string `json:"agent_start_time"`
	BootTime                 string `json:"boot_time"`
	RPMPatchInstallationTime string `json:"rpm_patch_installation_time"`
	IscsiIqn                 string `json:"iscsi_iqn"`
}

package entities

import "time"

import "encoding/json"

func UnmarshalPowerEntity(data []byte) (PowerEntity, error) {
	var r PowerEntity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PowerEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PowerEntity struct {
	OdataContext string        `json:"@odata.context"`
	OdataID      string        `json:"@odata.id"`
	OdataType    string        `json:"@odata.type"`
	Average      int64         `json:"Average"`
	ID           string        `json:"Id"`
	Maximum      int64         `json:"Maximum"`
	Minimum      int64         `json:"Minimum"`
	Name         string        `json:"Name"`
	PowerDetail  []PowerDetail `json:"PowerDetail"`
	Samples      int64         `json:"Samples"`
}

type PowerDetail struct {
	AmbTemp      int64     `json:"AmbTemp"`
	Average      int64     `json:"Average"`
	Cap          int64     `json:"Cap"`
	CPUAvgFreq   int64     `json:"CpuAvgFreq"`
	CPUCapLim    int64     `json:"CpuCapLim"`
	CPUMax       int64     `json:"CpuMax"`
	CPUPwrSavLim int64     `json:"CpuPwrSavLim"`
	CPUUtil      int64     `json:"CpuUtil"`
	Minimum      int64     `json:"Minimum"`
	Peak         int64     `json:"Peak"`
	PRMode       PRMode    `json:"PrMode"`
	PunCap       bool      `json:"PunCap"`
	Time         time.Time `json:"Time"`
	UnachCap     bool      `json:"UnachCap"`
}

type PRMode string

const (
	Dyn PRMode = "dyn"
)

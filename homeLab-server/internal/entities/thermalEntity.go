package entities

import "encoding/json"

func UnmarshalThermal(data []byte) (Thermal, error) {
	var r Thermal
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Thermal) Thermal() ([]byte, error) {
	return json.Marshal(r)
}

type Thermal struct {
	OdataContext string        `json:"@odata.context"`
	OdataID      string        `json:"@odata.id"`
	OdataType    string        `json:"@odata.type"`
	Fans         []Fan         `json:"Fans"`
	ID           string        `json:"Id"`
	Name         string        `json:"Name"`
	Temperatures []Temperature `json:"Temperatures"`
	Type         string        `json:"Type"`
	Links        Links         `json:"links"`
}

type Fan struct {
	CurrentReading int64  `json:"CurrentReading"`
	FanName        string `json:"FanName"`
	OEM            FanOEM `json:"Oem"`
	Status         Status `json:"Status"`
	Units          string `json:"Units"`
}

type FanOEM struct {
	HP PurpleHP `json:"Hp"`
}

type PurpleHP struct {
	OdataType string `json:"@odata.type"`
	Location  string `json:"Location"`
	Type      string `json:"Type"`
}

type Status struct {
	Health *Health `json:"Health,omitempty"`
	State  State   `json:"State"`
}

type Links struct {
	Self Self `json:"self"`
}

type Self struct {
	Href string `json:"href"`
}

type Temperature struct {
	CurrentReading         int64           `json:"CurrentReading"`
	Name                   string          `json:"Name"`
	Number                 int64           `json:"Number"`
	OEM                    TemperatureOEM  `json:"Oem"`
	PhysicalContext        PhysicalContext `json:"PhysicalContext"`
	ReadingCelsius         int64           `json:"ReadingCelsius"`
	Status                 Status          `json:"Status"`
	Units                  Units           `json:"Units"`
	UpperThresholdCritical int64           `json:"UpperThresholdCritical"`
	UpperThresholdFatal    int64           `json:"UpperThresholdFatal"`
	UpperThresholdUser     *int64          `json:"UpperThresholdUser,omitempty"`
}

type TemperatureOEM struct {
	HP FluffyHP `json:"Hp"`
}

type FluffyHP struct {
	OdataType   OdataType `json:"@odata.type"`
	LocationXmm int64     `json:"LocationXmm"`
	LocationYmm int64     `json:"LocationYmm"`
	Type        Type      `json:"Type"`
}

type Health string

const (
	Ok Health = "OK"
)

type State string

const (
	Absent  State = "Absent"
	Enabled State = "Enabled"
)

type OdataType string

const (
	HPSeaOfSensors100HPSeaOfSensors OdataType = "#HpSeaOfSensors.1.0.0.HpSeaOfSensors"
)

type Type string

const (
	HPSeaOfSensors100 Type = "HpSeaOfSensors.1.0.0"
)

type PhysicalContext string

const (
	CPU         PhysicalContext = "CPU"
	Intake      PhysicalContext = "Intake"
	PowerSupply PhysicalContext = "PowerSupply"
	SystemBoard PhysicalContext = "SystemBoard"
)

type Units string

const (
	Celsius Units = "Celsius"
)

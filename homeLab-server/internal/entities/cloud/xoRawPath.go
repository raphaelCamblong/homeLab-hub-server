package entities

import "encoding/json"

type XoRawPathEntity []string

func UnmarshalXoRawPathEntity(data []byte) (XoRawPathEntity, error) {
	var r XoRawPathEntity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *XoRawPathEntity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

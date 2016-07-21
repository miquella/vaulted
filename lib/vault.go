package vaulted

import (
	"encoding/base64"
)

type Vault struct {
	Vars map[string]string `json:"vars"`
}

type VaultFile struct {
	Key *VaultKey `json:"key"`

	Method     string  `json:"method"`
	Details    Details `json:"details,omitempty"`
	Ciphertext []byte  `json:"ciphertext"`
}

type Details map[string]interface{}

func (d Details) Clone() Details {
	newKeyDetails := make(Details)
	for k, v := range d {
		newKeyDetails[k] = v
	}
	return newKeyDetails
}

func (d Details) Int(name string) int {
	if v, ok := d[name].(int); ok {
		return v
	}
	if v, ok := d[name].(int64); ok {
		return int(v)
	}
	if v, ok := d[name].(float64); ok {
		return int(v)
	}
	return 0
}

func (d Details) SetInt(name string, value int) {
	d[name] = value
}

func (d Details) String(name string) string {
	v, _ := d[name].(string)
	return v
}

func (d Details) SetString(name string, value string) {
	d[name] = value
}

func (d Details) Bytes(name string) []byte {
	b, err := base64.StdEncoding.DecodeString(d.String(name))
	if err != nil {
		return nil
	}
	return b
}

func (d Details) SetBytes(name string, value []byte) {
	d[name] = base64.StdEncoding.EncodeToString(value)
}

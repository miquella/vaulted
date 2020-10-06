package vaulted

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/miquella/xdg"
	"golang.org/x/crypto/pbkdf2"
)

const (
	BaseIterations          = 1 << 17
	AdditionIterationsRange = 1 << 18
)

type VaultFile struct {
	Key *VaultKey `json:"key"`

	Method     string  `json:"method"`
	Details    Details `json:"details,omitempty"`
	Ciphertext []byte  `json:"ciphertext"`
}

func readVaultFile(name string) (*VaultFile, error) {
	existing := xdg.DATA.Find(filepath.Join("vaulted", name))
	if len(existing) == 0 {
		return nil, os.ErrNotExist
	}

	f, err := os.Open(existing[0])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	vf := VaultFile{}
	err = d.Decode(&vf)
	if err != nil {
		return nil, err
	}

	return &vf, nil
}

func writeVaultFile(name string, vaultFile *VaultFile) error {
	filename := xdg.DATA_HOME.Join(filepath.Join("vaulted", name))
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	err = e.Encode(vaultFile)
	if err != nil {
		return err
	}

	removeSessionCache(name)

	return nil
}

type VaultKey struct {
	Method  string  `json:"method"`
	Details Details `json:"details"`
}

func newVaultKey(previous *VaultKey) *VaultKey {
	var method string
	var details Details

	// Copy previous key details, if present
	if previous != nil {
		method = previous.Method
		details = previous.Details.Clone()
	} else {
		method = "pbkdf2-sha512"
		details = make(Details)

		iterations := BaseIterations
		r, err := rand.Int(rand.Reader, big.NewInt(AdditionIterationsRange))
		if err == nil {
			iterations += int(r.Int64())
		}

		details.SetInt("iterations", iterations)
	}

	// Generate new salt
	switch method {
	case "pbkdf2-sha512":
		salt := make([]byte, 32)
		_, err := rand.Read(salt)
		if err != nil {
			return nil
		}
		details.SetBytes("salt", salt)
	}

	return &VaultKey{
		Method:  method,
		Details: details,
	}
}

func (vk *VaultKey) key(password string, keyLength int) ([]byte, error) {
	switch vk.Method {
	case "pbkdf2-sha512":
		iterations := vk.Details.Int("iterations")
		salt := vk.Details.Bytes("salt")
		if iterations == 0 || len(salt) == 0 {
			return nil, ErrInvalidKeyConfig
		}
		return pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha512.New), nil
	}

	return nil, fmt.Errorf("Invalid key derivation method: %s", vk.Method)
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

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"

	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

func CaptureStdout(f func()) []byte {
	// Save/restore stdout
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()

	// Capture stdout
	r, w, _ := os.Pipe()
	defer w.Close()
	os.Stdout = w

	captured := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		captured <- buf.Bytes()
		close(captured)
	}()

	f()

	w.Close()
	return <-captured
}

func WriteStdin(b []byte, f func()) {
	// Save/restore stdin
	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()

	// Write to stdin
	r, w, _ := os.Pipe()
	defer r.Close()
	os.Stdin = r
	go func() {
		vr := bytes.NewReader(b)
		io.Copy(w, vr)
		w.Close()
	}()

	f()
}

func NewTestStore() *TestStore {
	return &TestStore{
		Passwords: make(map[string]string),
		Vaults:    make(map[string]*vaulted.Vault),
		Sessions:  make(map[string]*vaulted.Session),
	}
}

type TestStore struct {
	Passwords map[string]string
	Vaults    map[string]*vaulted.Vault
	Sessions  map[string]*vaulted.Session

	LegacyPassword     string
	LegacyEnvironments map[string]legacy.Environment
}

func (ts TestStore) GetPassword(operation vaulted.Operation, name string) (string, error) {
	switch operation {
	case vaulted.OpenOperation:
		return "prompted open password", nil
	case vaulted.SealOperation:
		return "prompted seal password", nil
	default:
		return "", errors.New("Unknown operation")
	}
}

func (ts TestStore) GetMFAToken(name string) (string, error) {
	return "123456", nil
}

func (ts TestStore) Steward() vaulted.Steward {
	return ts
}

func (ts TestStore) VaultExists(name string) bool {
	_, exists := ts.Vaults[name]
	return exists
}

func (ts TestStore) ListVaults() ([]string, error) {
	var vaults []string
	for name := range ts.Vaults {
		vaults = append(vaults, name)
	}
	return vaults, nil
}

func (ts TestStore) SealVault(vault *vaulted.Vault, name string) error {
	return ts.SealVaultWithPassword(vault, name, "sealed")
}

func (ts TestStore) SealVaultWithPassword(vault *vaulted.Vault, name, password string) error {
	ts.Passwords[name] = password
	ts.Vaults[name] = cloneVault(vault)

	return nil
}

func (ts TestStore) OpenVault(name string) (*vaulted.Vault, string, error) {
	return ts.OpenVaultWithPassword(name, "prompted password")
}

func (ts TestStore) OpenVaultWithPassword(name, password string) (*vaulted.Vault, string, error) {
	if !ts.VaultExists(name) {
		return nil, "", os.ErrNotExist
	}

	return cloneVault(ts.Vaults[name]), ts.Passwords[name], nil
}

func (ts TestStore) RemoveVault(name string) error {
	if !ts.VaultExists(name) {
		return os.ErrNotExist
	}

	delete(ts.Passwords, name)
	delete(ts.Vaults, name)

	return nil
}

func (ts TestStore) GetSession(name string) (*vaulted.Session, string, error) {
	if !ts.VaultExists(name) {
		return nil, "", os.ErrNotExist
	}

	s := &vaulted.Session{
		SessionVersion: vaulted.SessionVersion,

		Expiration: time.Unix(1136239445, 0),

		Vars:    make(map[string]string),
		SSHKeys: make(map[string]string),
	}
	if _, exists := ts.Sessions[name]; exists {
		cachedSession := ts.Sessions[name]

		s.Name = cachedSession.Name
		s.Expiration = cachedSession.Expiration

		creds := *cachedSession.AWSCreds
		s.AWSCreds = &creds

		for key, value := range cachedSession.Vars {
			s.Vars[key] = value
		}

		for key, value := range cachedSession.SSHKeys {
			s.SSHKeys[key] = value
		}
	} else {
		vault := ts.Vaults[name]

		s.Name = name

		for key, value := range vault.Vars {
			s.Vars[key] = value
		}

		for key, value := range vault.SSHKeys {
			s.SSHKeys[key] = value
		}
	}

	return s, ts.Passwords[name], nil
}

func (ts TestStore) CreateSession(name string) (*vaulted.Session, string, error) {
	if !ts.VaultExists(name) {
		return nil, "", os.ErrNotExist
	}

	s := &vaulted.Session{
		SessionVersion: vaulted.SessionVersion,

		Expiration: time.Unix(1136239446, 0),

		Vars:    make(map[string]string),
		SSHKeys: make(map[string]string),
	}
	vault := ts.Vaults[name]

	s.Name = name

	for key, value := range vault.Vars {
		s.Vars[key] = value
	}

	for key, value := range vault.SSHKeys {
		s.SSHKeys[key] = value
	}
	return s, ts.Passwords[name], nil
}

func (ts TestStore) OpenLegacyVault() (environments map[string]legacy.Environment, password string, err error) {
	return ts.LegacyEnvironments, ts.LegacyPassword, nil
}

func cloneVault(vault *vaulted.Vault) *vaulted.Vault {
	newVault := &vaulted.Vault{
		Vars:    make(map[string]string),
		SSHKeys: make(map[string]string),
	}

	if vault.AWSKey != nil {
		newVault.AWSKey = &vaulted.AWSKey{
			AWSCredentials: vaulted.AWSCredentials{
				ID:     vault.AWSKey.ID,
				Secret: vault.AWSKey.Secret,
			},
			MFA:  vault.AWSKey.MFA,
			Role: vault.AWSKey.Role,
		}
	}

	for key, value := range vault.Vars {
		newVault.Vars[key] = value
	}

	for key, value := range vault.SSHKeys {
		newVault.SSHKeys[key] = value
	}

	return newVault
}

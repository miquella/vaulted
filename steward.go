package main

import (
	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

type Steward interface {
	VaultExists(name string) bool
	ListVaults() ([]string, error)
	SealVault(name string, password *string, vault *vaulted.Vault) error
	OpenVault(name string, password *string) (string, *vaulted.Vault, error)
	RemoveVault(name string) error
	GetSession(name string, password *string) (string, *vaulted.Session, error)

	OpenLegacyVault() (password string, environments map[string]legacy.Environment, err error)
}

type Command interface {
	Run(steward Steward) error
}

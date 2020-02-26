package main

import (
	"github.com/miquella/vaulted/v3/lib"
)

type Copy struct {
	OldVaultName string
	NewVaultName string
}

func (c *Copy) Run(store vaulted.Store) error {
	vault, _, err := store.OpenVault(c.OldVaultName)
	if err != nil {
		return err
	}

	err = store.SealVault(vault, c.NewVaultName)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"encoding/json"
	"os"

	"github.com/miquella/vaulted/v3/lib"
)

type Dump struct {
	VaultName string
}

func (d *Dump) Run(store vaulted.Store) error {
	vault, _, err := store.OpenVault(d.VaultName)
	if err != nil {
		return err
	}

	jvault, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return err
	}

	for len(jvault) > 0 {
		n, err := os.Stdout.Write(jvault)
		if err != nil {
			return err
		}

		jvault = jvault[n:]
	}

	return nil
}

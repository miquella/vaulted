package main

import (
	"encoding/json"
	"os"
)

type Dump struct {
	VaultName string
}

func (d *Dump) Run(steward Steward) error {
	_, vault, err := steward.OpenVault(d.VaultName, nil)
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

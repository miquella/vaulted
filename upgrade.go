package main

import (
	"errors"
	"fmt"

	"github.com/miquella/vaulted/lib"
)

var (
	ErrUpgradeFailed = errors.New("Upgrade failed")
)

type Upgrade struct{}

func (u *Upgrade) Run(steward Steward) error {
	password, environments, err := steward.OpenLegacyVault()
	if err != nil {
		return err
	}

	// collect the current list of vaults (so we don't overwrite any)
	vaults, _ := steward.ListVaults()
	existingVaults := map[string]bool{}
	for _, name := range vaults {
		existingVaults[name] = true
	}

	failed := 0
	for name, env := range environments {
		if existingVaults[name] {
			fmt.Printf("%s: skipped (vault already exists)\n", name)
			continue
		}

		vault := vaulted.Vault{
			Vars: env.Vars,
		}
		err = steward.SealVault(name, &password, &vault)
		if err != nil {
			failed++
			fmt.Printf("%s: %v\n", name, err)
		} else {
			fmt.Printf("%s: upgraded\n", name)
		}
	}

	if failed > 0 {
		return ErrorWithExitCode{ErrUpgradeFailed, failed}
	}

	return nil
}

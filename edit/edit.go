package edit

import (
	"errors"
	"fmt"

	"github.com/miquella/vaulted/v3/lib"
	"github.com/miquella/vaulted/v3/menu"
)

var (
	// ErrExists error for when the vault exists already on "create"
	ErrExists = errors.New("vault exists. Use `vaulted edit` to edit existing vaults")
)

type Edit struct {
	New       bool
	VaultName string
}

func (e *Edit) Run(store vaulted.Store) error {
	var password string
	var vault *vaulted.Vault
	var err error

	if e.New {
		if store.VaultExists(e.VaultName) {
			return ErrExists
		}

		fmt.Printf("Creating new vault '%s'...\n", e.VaultName)
		vault = &vaulted.Vault{}

		importCredsMenu := menu.ImportCredentialsMenu{}
		err := importCredsMenu.Handler()
		if err != nil {
			return err
		}

		creds := importCredsMenu.Credentials
		if creds != nil {
			vault.AWSKey = &vaulted.AWSKey{
				AWSCredentials: *creds,
			}

			detectMFAMenu := menu.DetectMFAMenu{Menu: &menu.Menu{Vault: vault}}
			_ = detectMFAMenu.Handler()
		}

	} else {
		vault, password, err = store.OpenVault(e.VaultName)
		if err != nil {
			return err
		}
	}

	err = e.edit(e.VaultName, vault)
	if err != nil {
		return err
	}

	if e.New {
		err = store.SealVault(vault, e.VaultName)
	} else {
		err = store.SealVaultWithPassword(vault, e.VaultName, password)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Vault '%s' successfully saved!\n", e.VaultName)

	return nil
}

func (e *Edit) edit(name string, v *vaulted.Vault) error {
	mainMenu := &menu.MainMenu{
		Menu: menu.Menu{
			Vault: v,
		},
		VaultName: e.VaultName,
	}

	return mainMenu.Handler()
}

package main

import (
	"fmt"
	"os"

	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

type TTYSteward struct{}

func (*TTYSteward) VaultExists(name string) bool {
	return vaulted.VaultExists(name)
}

func (*TTYSteward) ListVaults() ([]string, error) {
	return vaulted.ListVaults()
}

func (*TTYSteward) SealVault(name string, password *string, vault *vaulted.Vault) error {
	if password == nil {
		envPassword := os.Getenv("VAULTED_NEW_PASSWORD")
		if envPassword != "" {
			password = &envPassword
		} else {
			for {
				newPassword, err := ask.HiddenAsk("New Password: ")
				if err != nil {
					return err
				}
				confirmPassword, err := ask.HiddenAsk("Confirm Password: ")
				if err != nil {
					return err
				}

				if newPassword == confirmPassword {
					password = &newPassword
					break
				} else {
					ask.Print("Passwords do not match.\n")
				}
			}
		}
	}

	return vaulted.SealVault(*password, name, vault)
}

func (*TTYSteward) OpenVault(name string, password *string) (string, *vaulted.Vault, error) {
	if !vaulted.VaultExists(name) {
		return "", nil, os.ErrNotExist
	}

	if password == nil && os.Getenv("VAULTED_PASSWORD") != "" {
		envPassword := os.Getenv("VAULTED_PASSWORD")
		password = &envPassword
	}

	var vault *vaulted.Vault
	var err error
	if password != nil {
		vault, err = vaulted.OpenVault(*password, name)
	} else {
		for i := 0; i < 3; i++ {
			var requestedPassword string
			requestedPassword, err = ask.HiddenAsk("Password: ")
			if err != nil {
				break
			}

			vault, err = vaulted.OpenVault(requestedPassword, name)
			if err != vaulted.ErrInvalidPassword {
				password = &requestedPassword
				break
			}
		}
	}

	if err != nil {
		return "", nil, err
	}

	return *password, vault, nil
}

func (*TTYSteward) RemoveVault(name string) error {
	return vaulted.RemoveVault(name)
}

func (*TTYSteward) GetEnvironment(name string, password *string) (string, *vaulted.Environment, error) {
	if !vaulted.VaultExists(name) {
		return "", nil, os.ErrNotExist
	}

	if password == nil && os.Getenv("VAULTED_PASSWORD") != "" {
		envPassword := os.Getenv("VAULTED_PASSWORD")
		password = &envPassword
	}

	var env *vaulted.Environment
	var err error
	if password != nil {
		env, err = vaulted.GetEnvironment(*password, name)
	} else {
		for i := 0; i < 3; i++ {
			var requestedPassword string
			requestedPassword, err = ask.HiddenAsk("Password: ")
			if err != nil {
				break
			}

			env, err = vaulted.GetEnvironment(requestedPassword, name)
			if err != vaulted.ErrInvalidPassword {
				password = &requestedPassword
				break
			}
		}
	}

	if err != nil {
		return "", nil, err
	}

	return *password, env, nil
}

func (*TTYSteward) OpenLegacyVault() (password string, environments map[string]legacy.Environment, err error) {
	legacyVault, err := legacy.ReadVault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	password = os.Getenv("VAULTED_PASSWORD")
	if password != "" {
		environments, err = legacyVault.DecryptEnvironments(password)
	} else {
		for i := 0; i < 3; i++ {
			password, err = ask.HiddenAsk("Legacy Password: ")
			if err != nil {
				break
			}

			environments, err = legacyVault.DecryptEnvironments(password)
			if err != legacy.ErrInvalidPassword {
				break
			}
		}
	}
	return
}

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

var (
	ErrFileNotExist      = ErrorWithExitCode{os.ErrNotExist, EX_USAGE_ERROR}
	ErrNoPasswordEntered = ErrorWithExitCode{errors.New("Could not get password"), EX_UNAVAILABLE}

	vaultedErrMap = map[error]ErrorWithExitCode{
		vaulted.ErrInvalidPassword:         ErrorWithExitCode{vaulted.ErrInvalidPassword, EX_TEMPORARY_ERROR},
		vaulted.ErrInvalidKeyConfig:        ErrorWithExitCode{vaulted.ErrInvalidKeyConfig, EX_DATA_ERROR},
		vaulted.ErrInvalidEncryptionConfig: ErrorWithExitCode{vaulted.ErrInvalidEncryptionConfig, EX_DATA_ERROR},
	}
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
				newPassword, err := getPassword("New Password: ")
				if err != nil {
					return err
				}
				confirmPassword, err := getPassword("Confirm Password: ")
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

	return vaulted.SealVault(name, *password, vault)
}

func (*TTYSteward) OpenVault(name string, password *string) (string, *vaulted.Vault, error) {
	if !vaulted.VaultExists(name) {
		return "", nil, ErrFileNotExist
	}

	if password == nil && os.Getenv("VAULTED_PASSWORD") != "" {
		envPassword := os.Getenv("VAULTED_PASSWORD")
		password = &envPassword
	}

	var vault *vaulted.Vault
	var err error
	if password != nil {
		vault, err = vaulted.OpenVault(name, *password)
	} else {
		for i := 0; i < 3; i++ {
			var requestedPassword string
			requestedPassword, err = getPassword("Password: ")
			if err != nil {
				break
			}

			vault, err = vaulted.OpenVault(name, requestedPassword)
			if err != vaulted.ErrInvalidPassword {
				password = &requestedPassword
				break
			}
		}
	}

	if err != nil {
		if _, ok := vaultedErrMap[err]; ok {
			return "", nil, vaultedErrMap[err]
		} else {
			return "", nil, err
		}
	}

	return *password, vault, nil
}

func (*TTYSteward) RemoveVault(name string) error {
	return vaulted.RemoveVault(name)
}

func (*TTYSteward) GetSession(name string, password *string) (string, *vaulted.Session, error) {
	if !vaulted.VaultExists(name) {
		return "", nil, ErrFileNotExist
	}

	if password == nil && os.Getenv("VAULTED_PASSWORD") != "" {
		envPassword := os.Getenv("VAULTED_PASSWORD")
		password = &envPassword
	}

	var session *vaulted.Session
	var err error
	if password != nil {
		session, err = vaulted.GetSession(name, *password)
	} else {
		for i := 0; i < 3; i++ {
			var requestedPassword string
			requestedPassword, err = getPassword("Password: ")
			if err != nil {
				break
			}

			session, err = vaulted.GetSession(name, requestedPassword)
			if err != vaulted.ErrInvalidPassword {
				password = &requestedPassword
				break
			}
		}
	}

	if err != nil {
		if _, ok := vaultedErrMap[err]; ok {
			return "", nil, vaultedErrMap[err]
		} else {
			return "", nil, err
		}
	}

	return *password, session, nil
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
			password, err = getPassword("Legacy Password: ")
			if err != nil {
				break
			}

			environments, err = legacyVault.DecryptEnvironments(password)
			if err != vaulted.ErrInvalidPassword {
				break
			}
		}
	}
	return
}

func getPassword(prompt string) (string, error) {
	if os.Getenv("VAULTED_ASKPASS") != "" {
		cmd := exec.Command(os.Getenv("VAULTED_ASKPASS"), prompt)
		output, err := cmd.Output()
		if err != nil {
			return "", ErrNoPasswordEntered
		}

		return strings.Trim(string(output), "\r\n"), nil

	} else {
		return ask.HiddenAsk(prompt)
	}
}

package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

type TTYSteward struct{}

func (t *TTYSteward) GetMaxOpenTries() int {
	if os.Getenv("VAULTED_PASSWORD") != "" {
		return 1
	}

	// we'll try tty prompt & askpass up to 3 times
	return 3
}

func (t *TTYSteward) GetPassword(operation vaulted.Operation, name string) (string, error) {
	switch operation {
	case vaulted.SealOperation:
		if password := os.Getenv("VAULTED_NEW_PASSWORD"); password != "" {
			return password, nil
		}

		for {
			password, err := t.getPassword("New password: ")
			if err != nil {
				return "", err
			}

			confirm, err := t.getPassword("Confirm Password: ")
			if err != nil {
				return "", err
			}

			if password == confirm {
				return password, nil
			}

			ask.Print("Passwords do not match.\n")
		}

	case legacy.LegacyOperation:
		return t.getPassword("Legacy Password: ")

	default:
		if password := os.Getenv("VAULTED_PASSWORD"); password != "" {
			return password, nil
		}

		return t.getPassword("Password: ")
	}
}

func (t *TTYSteward) getPassword(prompt string) (string, error) {
	if askPass := os.Getenv("VAULTED_ASKPASS"); askPass != "" {
		cmd := exec.Command(askPass, prompt)
		output, err := cmd.Output()
		if err != nil {
			return "", ErrNoPasswordEntered
		}

		return strings.Trim(string(output), "\r\n"), nil
	}

	return ask.HiddenAsk(prompt)
}

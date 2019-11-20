package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/miquella/ask"

	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

const (
	PASSWORD_TYPE_PASSWORD        = "password"
	PASSWORD_TYPE_LEGACY_PASSWORD = "legacypassword"
	PASSWORD_TYPE_MFATOKEN        = "mfatoken"

	PASSWORD_REASON_NEW     = "new"
	PASSWORD_REASON_NOMATCH = "nomatch"
	PASSWORD_REASON_CONFIRM = "confirm"
)

func NewSteward() vaulted.Steward {
	if askpass, present := os.LookupEnv("VAULTED_ASKPASS"); present {
		return &AskPassSteward{
			Command: askpass,
		}
	} else {
		return &TTYSteward{}
	}
}

type AskPassSteward struct {
	Command string
}

func (t *AskPassSteward) GetMaxOpenTries() int {
	if _, present := os.LookupEnv("VAULTED_PASSWORD"); present {
		return 1
	}

	// we'll try askpass up to 3 times
	return 3
}

func (t *AskPassSteward) GetPassword(operation vaulted.Operation, name string) (string, error) {
	// environment variables take precedence
	switch operation {
	case vaulted.SealOperation:
		if password, present := os.LookupEnv("VAULTED_NEW_PASSWORD"); present {
			return password, nil
		}

	case legacy.LegacyOperation:

	default:
		if password, present := os.LookupEnv("VAULTED_PASSWORD"); present {
			return password, nil
		}
	}

	// askpass prompt
	switch operation {
	case vaulted.SealOperation:
		for firstTry := false; ; firstTry = true {
			var passwordreason string
			var prompt string
			if firstTry {
				passwordreason = PASSWORD_REASON_NEW
				prompt = fmt.Sprintf("'%s' new password: ", name)
			} else {
				passwordreason = PASSWORD_REASON_NOMATCH
				prompt = fmt.Sprintf("'%s' new password (passwords didn't match): ", name)
			}
			password, err := t.askpass(name, PASSWORD_TYPE_PASSWORD, passwordreason, prompt)
			if err != nil {
				return "", err
			}

			confirm, err := t.askpass(name, PASSWORD_TYPE_PASSWORD, PASSWORD_REASON_CONFIRM, fmt.Sprintf("'%s' confirm password: ", name))
			if err != nil {
				return "", err
			}

			if password == confirm {
				return password, nil
			}
		}

	case legacy.LegacyOperation:
		return t.askpass(name, PASSWORD_TYPE_LEGACY_PASSWORD, "", "Legacy Password: ")

	default:
		return t.askpass(name, PASSWORD_TYPE_PASSWORD, "", fmt.Sprintf("'%s' password: ", name))
	}
}

func (t *AskPassSteward) GetMFAToken(name string) (string, error) {
	return t.askpass(name, PASSWORD_TYPE_MFATOKEN, "", fmt.Sprintf("'%s' MFA token: ", name))
}

func (t *AskPassSteward) askpass(name string, passwordtype string, reason string, prompt string) (string, error) {
	cmd := exec.Command(t.Command, prompt)
	cmd.Env = append(os.Environ(),
		"VAULTED_ENV="+name,
		"VAULTED_PASSWORD_TYPE="+passwordtype,
		"VAULTED_PASSWORD_REASON="+reason,
	)
	output, err := cmd.Output()
	if err != nil {
		return "", ErrNoPasswordEntered
	}

	return strings.Trim(string(output), "\r\n"), nil
}

type TTYSteward struct{}

func (t *TTYSteward) GetMaxOpenTries() int {
	if _, present := os.LookupEnv("VAULTED_PASSWORD"); present {
		return 1
	}

	// we'll try tty prompt up to 3 times
	return 3
}

func (t *TTYSteward) GetPassword(operation vaulted.Operation, name string) (string, error) {
	// environment variables take precedence
	switch operation {
	case vaulted.SealOperation:
		if password, present := os.LookupEnv("VAULTED_NEW_PASSWORD"); present {
			return password, nil
		}

	case legacy.LegacyOperation:

	default:
		if password, present := os.LookupEnv("VAULTED_PASSWORD"); present {
			return password, nil
		}
	}

	// tty prompt
	switch operation {
	case vaulted.SealOperation:
		ask.Print(fmt.Sprintf("Vault '%s'\n", name))
		for {
			password, err := ask.HiddenAsk("   New password: ")
			if err != nil {
				return "", err
			}

			confirm, err := ask.HiddenAsk("   Confirm password: ")
			if err != nil {
				return "", err
			}

			if password == confirm {
				return password, nil
			}

			ask.Print("Passwords do not match.\n\n")
		}

	case legacy.LegacyOperation:
		return ask.HiddenAsk("Legacy Password: ")

	default:
		ask.Print(fmt.Sprintf("Vault '%s'\n", name))
		return ask.HiddenAsk("   Password: ")
	}
}

var (
	mfaTokenValidation = regexp.MustCompile(`^\d{6}$`)
)

func (t *TTYSteward) GetMFAToken(name string) (string, error) {
	for attempts := 0; attempts < 3; attempts++ {
		token, err := ask.Ask("   MFA token: ")
		if err != nil {
			return "", err
		}

		token = strings.TrimSpace(token)
		if mfaTokenValidation.MatchString(token) {
			return token, nil
		}

		ask.Print("Invalid MFA token.\n")
	}

	return "", ErrNoMFATokenEntered
}

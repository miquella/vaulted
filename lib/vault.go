package vaulted

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/miquella/ask"
)

const (
	DefaultSessionName = "VaultedSession"
)

var STSDurationDefault = time.Hour

var (
	ErrInvalidCommand = errors.New("Invalid command")
	ErrNoTokenEntered = errors.New("Could not get MFA code")
)

type Vault struct {
	Duration time.Duration     `json:"duration,omitempty"`
	AWSKey   *AWSKey           `json:"aws_key,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
	SSHKeys  map[string]string `json:"ssh_keys,omitempty"`
}

func (v *Vault) CreateSession(name string) (*Session, error) {
	var duration time.Duration
	if v.Duration == 0 {
		duration = STSDurationDefault
	} else {
		duration = v.Duration
	}

	s := &Session{
		Name:       name,
		Vars:       make(map[string]string),
		Expiration: time.Now().Add(duration).Truncate(time.Second),
	}

	// copy the vault vars to the session
	for key, value := range v.Vars {
		s.Vars[key] = value
	}

	// copy the vault ssh keys to the session
	if len(v.SSHKeys) > 0 {
		s.SSHKeys = make(map[string]string)
		for key, value := range v.SSHKeys {
			s.SSHKeys[key] = value
		}
	}

	// get aws creds
	if v.AWSKey != nil && v.AWSKey.ID != "" && v.AWSKey.Secret != "" {
		var err error
		s.AWSCreds, err = v.AWSKey.GetAWSCredentials(duration)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

type AWSKey struct {
	AWSCredentials
	MFA                     string `json:"mfa,omitempty"`
	Role                    string `json:"role,omitempty"`
	ForgoTempCredGeneration bool   `json:"forgoTempCredGeneration"`
}

func (k *AWSKey) GetAWSCredentials(duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	if k.MFA == "" {
		return k.AWSCredentials.GetSessionToken(duration)
	}

	tokenCode, err := getTokenCode()
	if err != nil {
		return nil, err
	}

	return k.AWSCredentials.GetSessionTokenWithMFA(k.MFA, tokenCode, duration)
}

func getTokenCode() (string, error) {
	prompt := "Enter your MFA code: "
	if os.Getenv("VAULTED_ASKPASS") != "" {
		cmd := exec.Command(os.Getenv("VAULTED_ASKPASS"), prompt)
		output, err := cmd.Output()
		if err != nil {
			return "", ErrNoTokenEntered
		}

		return strings.TrimSpace(string(output)), nil

	} else {
		tokenCode, err := ask.Ask(prompt)
		return strings.TrimSpace(tokenCode), err
	}
}

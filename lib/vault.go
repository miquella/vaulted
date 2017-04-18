package vaulted

import (
	"errors"
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
)

type Vault struct {
	Duration time.Duration     `json:"duration,omitempty"`
	AWSKey   *AWSKey           `json:"aws_key,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
	SSHKeys  map[string]string `json:"ssh_keys,omitempty"`
}

func (v *Vault) CreateEnvironment(extraVars map[string]string) (*Environment, error) {
	var duration time.Duration
	if v.Duration == 0 {
		duration = STSDurationDefault
	} else {
		duration = v.Duration
	}

	e := &Environment{
		Vars:       make(map[string]string),
		Expiration: time.Now().Add(duration),
	}

	// copy the vault vars to the environment
	for key, value := range v.Vars {
		e.Vars[key] = value
	}
	for key, value := range extraVars {
		e.Vars[key] = value
	}

	// copy the vault ssh keys to the environment
	if len(v.SSHKeys) > 0 {
		e.SSHKeys = make(map[string]string)
		for key, value := range v.SSHKeys {
			e.SSHKeys[key] = value
		}
	}

	// get aws creds
	if v.AWSKey != nil && v.AWSKey.ID != "" && v.AWSKey.Secret != "" {
		var err error
		e.AWSCreds, err = v.AWSKey.GetAWSCredentials(duration)
		if err != nil {
			return nil, err
		}

		if v.AWSKey.Role != "" {
			e.AWSCreds, err = e.AWSCreds.AssumeRole(v.AWSKey.Role, duration)
			if err != nil {
				return nil, err
			}
		}
	}

	return e, nil
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
	tokenCode, err := ask.Ask("Enter your MFA code: ")
	if err != nil {
		return "", err
	}
	tokenCode = strings.TrimSpace(tokenCode)
	return tokenCode, nil
}

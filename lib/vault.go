package vaulted

import (
	"errors"
	"time"
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
	return v.createSession(name, func(duration time.Duration) (*AWSCredentials, error) {
		return v.AWSKey.GetAWSCredentials(duration)
	})
}

func (v *Vault) CreateSessionWithMFA(name, mfaToken string) (*Session, error) {
	return v.createSession(name, func(duration time.Duration) (*AWSCredentials, error) {
		return v.AWSKey.GetAWSCredentialsWithMFA(mfaToken, duration)
	})
}

func (v *Vault) createSession(name string, credsFunc func(duration time.Duration) (*AWSCredentials, error)) (*Session, error) {
	var duration time.Duration
	if v.Duration == 0 {
		duration = STSDurationDefault
	} else {
		duration = v.Duration
	}

	s := &Session{
		Name: name,
		Vars: make(map[string]string),
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

	if v.AWSKey.Valid() {
		var err error
		s.AWSCreds, err = credsFunc(duration)
		if err != nil {
			return nil, err
		}
	}

	// now that the session is generated, set the expiration
	s.Expiration = time.Now().Add(duration).Truncate(time.Second)

	return s, nil
}

type AWSKey struct {
	AWSCredentials
	MFA                     string `json:"mfa,omitempty"`
	Role                    string `json:"role,omitempty"`
	ForgoTempCredGeneration bool   `json:"forgoTempCredGeneration"`
}

func (k *AWSKey) Valid() bool {
	return k != nil && k.AWSCredentials.Valid()
}

func (k *AWSKey) RequiresMFA() bool {
	return k.Valid() && !k.ForgoTempCredGeneration && k.MFA != ""
}

func (k *AWSKey) GetAWSCredentials(duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	return k.AWSCredentials.GetSessionToken(duration)
}

func (k *AWSKey) GetAWSCredentialsWithMFA(mfaToken string, duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	return k.AWSCredentials.GetSessionTokenWithMFA(k.MFA, mfaToken, duration)
}

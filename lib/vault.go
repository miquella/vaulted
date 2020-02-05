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

type SSHOptions struct {
	DisableProxy    bool     `json:"disable_proxy"`
	GenerateRSAKey  bool     `json:"generate_rsa_key"`
	ValidPrincipals []string `json:"valid_principals,omitempty"`
	VaultSigningUrl string   `json:"vault_signing_url,omitempty"`
}

type Vault struct {
	Duration   time.Duration     `json:"duration,omitempty"`
	AWSKey     *AWSKey           `json:"aws_key,omitempty"`
	Vars       map[string]string `json:"vars,omitempty"`
	SSHKeys    map[string]string `json:"ssh_keys,omitempty"`
	SSHOptions *SSHOptions       `json:"ssh_options,omitempty"`
}

func (v *Vault) NewSession(name string) (*Session, error) {
	return v.newSession(name, func(duration time.Duration) (*AWSCredentials, error) {
		return v.AWSKey.GetAWSCredentials(duration)
	})
}

func (v *Vault) NewSessionWithMFA(name, mfaToken string) (*Session, error) {
	return v.newSession(name, func(duration time.Duration) (*AWSCredentials, error) {
		return v.AWSKey.GetAWSCredentialsWithMFA(mfaToken, duration)
	})
}

func (v *Vault) newSession(name string, credsFunc func(duration time.Duration) (*AWSCredentials, error)) (*Session, error) {
	var duration time.Duration
	if v.Duration == 0 {
		duration = STSDurationDefault
	} else {
		duration = v.Duration
	}

	var expiration *time.Time

	s := &Session{
		Name: name,

		SSHOptions: &SSHOptions{},
		Vars:       make(map[string]string),
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

	// copy the vault ssh options to the session
	if v.SSHOptions != nil {
		s.SSHOptions = v.SSHOptions
	}

	if v.AWSKey.Valid() {
		var err error
		s.AWSCreds, err = credsFunc(duration)
		if err != nil {
			return nil, err
		}
		s.Role = v.AWSKey.Role

		expiration = s.AWSCreds.Expiration
	}

	// now that the session is generated, set the expiration
	if expiration != nil {
		s.Expiration = *expiration
	} else {
		s.Expiration = time.Now().Add(duration).Truncate(time.Second)
	}

	return s, nil
}

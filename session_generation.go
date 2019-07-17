package main

import (
	"errors"
	"os"
	"time"

	"github.com/miquella/vaulted/lib"
)

var (
	ErrNoSessionIncompatibleWithAssume  = errors.New("--assume generates session credentials, it cannot be combined with --no-session")
	ErrNoSessionIncompatibleWithRefresh = errors.New("--refresh refreshes session credentials, it cannot be combined with --no-session")
	ErrNoSessionRequiresVaultName       = errors.New("A vault name must be specified when using --no-session")
)

type SessionOptions struct {
	VaultName string

	NoSession bool

	Refresh bool
	Role    string
}

func GetSessionWithOptions(store vaulted.Store, options *SessionOptions) (*vaulted.Session, error) {
	// Disabled session credentials
	if options.NoSession {
		if options.VaultName == "" {
			return nil, ErrNoSessionRequiresVaultName
		} else if options.Refresh {
			return nil, ErrNoSessionIncompatibleWithRefresh
		} else if options.Role != "" {
			return nil, ErrNoSessionIncompatibleWithAssume
		}

		return getVaultSessionWithNoSession(store, options)
	}

	var session *vaulted.Session
	var err error

	// Get a session
	if options.VaultName == "" {
		session, err = getDefaultSession(options)
	} else {
		session, err = getVaultSession(store, options)
	}
	if err != nil {
		return nil, err
	}

	// Assume any role specified
	if options.Role != "" {
		return session.AssumeRole(options.Role)
	}

	return session, nil
}

func getVaultSessionWithNoSession(store vaulted.Store, options *SessionOptions) (*vaulted.Session, error) {
	vault, _, err := store.OpenVault(options.VaultName)
	if err != nil {
		return nil, err
	}

	// Change the in-memory vault to forgo temp cred generation
	if vault.AWSKey != nil {
		vault.AWSKey.ForgoTempCredGeneration = true
	}

	// Skip assuming the vault's role

	return vault.NewSession(options.VaultName)
}

func getDefaultSession(options *SessionOptions) (*vaulted.Session, error) {
	// Create the vault
	vault := &vaulted.Vault{
		Duration: time.Hour,
	}

	updateVaultFromEnvAndOptions(vault, options)

	// Create the session
	return vault.NewSession(os.Getenv("VAULTED_ENV"))
}

func getVaultSession(store vaulted.Store, options *SessionOptions) (*vaulted.Session, error) {
	vault, password, err := store.OpenVault(options.VaultName)
	if err != nil {
		return nil, err
	}

	updateVaultFromEnvAndOptions(vault, options)

	// Create/get cached session
	var session *vaulted.Session
	if options.Refresh {
		session, err = store.CreateSession(vault, options.VaultName, password)
	} else {
		session, err = store.GetSession(vault, options.VaultName, password)
	}
	if err != nil {
		return nil, err
	}

	// Assume the session's role
	return session.AssumeSessionRole()
}

func updateVaultFromEnvAndOptions(vault *vaulted.Vault, options *SessionOptions) {
	// Calculate the region (lowest precedence to highest)
	region := os.Getenv("AWS_DEFAULT_REGION")
	if awsRegion := os.Getenv("AWS_REGION"); awsRegion != "" {
		region = awsRegion
	}
	if vault.AWSKey != nil {
		if vault.AWSKey.Region != nil && *vault.AWSKey.Region != "" {
			region = *vault.AWSKey.Region
		}
	}

	// Set the region
	if region != "" {
		if vault.AWSKey == nil {
			vault.AWSKey = &vaulted.AWSKey{}
		}

		vault.AWSKey.Region = &region
	}
}

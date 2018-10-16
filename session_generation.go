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

	Refresh      bool
	Role         string
	RoleDuration *time.Duration
}

func DefaultSession() *vaulted.Session {
	return &vaulted.Session{
		SessionVersion: vaulted.SessionVersion,

		Name:       os.Getenv("VAULTED_ENV"),
		Expiration: time.Now().Add(time.Hour).Truncate(time.Second),
	}
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

	var err error
	var session *vaulted.Session

	// Get a session
	if options.VaultName == "" {
		session = DefaultSession()
	} else {
		session, err = getVaultSession(store, options)
		if err != nil {
			return nil, err
		}
	}

	// Assume any role specified
	if options.Role != "" {
		return session.AssumeRole(options.Role, options.RoleDuration)
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

func getVaultSession(store vaulted.Store, options *SessionOptions) (*vaulted.Session, error) {
	var session *vaulted.Session
	var err error

	// Create/get cached session
	if options.Refresh {
		session, _, err = store.CreateSession(options.VaultName)
	} else {
		session, _, err = store.GetSession(options.VaultName)
	}
	if err != nil {
		return nil, err
	}

	// Assume the session's role
	if options.RoleDuration != nil {
		session.RoleDuration = options.RoleDuration
	}

	return session.AssumeSessionRole()
}

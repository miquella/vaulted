package vaulted

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"sort"
)

var (
	// SessionCacheVersion indicates the current version of the cache format.
	//
	// Any cache loaded that does not match this version is ignored. This
	// causes all caches written for previous versions to be invalidated.
	SessionCacheVersion = "3"
)

var (
	// ErrVaultSessionNotFound occurs when attempting to locate a vault session
	// in a SessionCache that isn't present.
	ErrVaultSessionNotFound = errors.New("Vault session not found")
)

// SessionCache stores sessions keyed based on the contents of the vault that
// spawned the session.
//
// See VaultSessionCacheKey for details on how the key is generated.
type SessionCache struct {
	SessionCacheVersion string              `json:"version"`
	Sessions            map[string]*Session `json:"sessions"`
}

// GetVaultSession retrieves a copy of a session in the cache.
//
// The retrieved session is keyed using the contents of the provided vault.
func (sc *SessionCache) GetVaultSession(vault *Vault) (*Session, error) {
	sessionKey := VaultSessionCacheKey(vault)
	if session, exists := sc.Sessions[sessionKey]; exists {
		return session.Clone(), nil
	}

	return nil, ErrVaultSessionNotFound
}

// PutVaultSession stores a copy of a session in the cache.
//
// The stored session is keyed using the contents of the provided vault.
func (sc *SessionCache) PutVaultSession(vault *Vault, session *Session) {
	if sc.Sessions == nil {
		sc.Sessions = make(map[string]*Session)
	}

	sessionKey := VaultSessionCacheKey(vault)
	sc.Sessions[sessionKey] = session.Clone()
}

// RemoveExpiredSessions removes sessions from the cache that have expired.
func (sc *SessionCache) RemoveExpiredSessions() {
	for key, session := range sc.Sessions {
		if session.Expired(NoTolerance) {
			delete(sc.Sessions, key)
		}
	}
}

// VaultSessionCacheKey computes a stable key based on the contents of a vault.
//
// The computed key is intended to be used for things such as a session cache.
func VaultSessionCacheKey(vault *Vault) string {
	// gather all of the key attributes
	keyAttributes := map[string]string{}

	if vault.AWSKey != nil {
		keyAttributes["aws_key_id"] = vault.AWSKey.ID
		keyAttributes["aws_key_secret"] = vault.AWSKey.Secret
		if vault.AWSKey.Region != nil {
			keyAttributes["aws_key_region"] = *vault.AWSKey.Region
		}

		keyAttributes["aws_key_mfa"] = vault.AWSKey.MFA
		keyAttributes["aws_key_role"] = vault.AWSKey.Role

		if vault.AWSKey.ForgoTempCredGeneration {
			keyAttributes["aws_key_sts"] = "false"
		} else {
			keyAttributes["aws_key_sts"] = "true"
		}
	}

	for key, value := range vault.Vars {
		keyAttributes["vars_"+key] = value
	}

	for key, value := range vault.SSHKeys {
		keyAttributes["ssh_key_"+key] = value
	}

	// we cannot compare the actual generated key, so instead we just
	// want to confirm that if its existence matches the current vault
	if vault.SSHOptions.GenerateRSAKey {
		keyAttributes["generated_key_exists"] = "true"
	} else {
		keyAttributes["generated_key_exists"] = "false"
	}

	// get a sorted list of the keys (that do not have blank values)
	var keys []string
	for key, value := range keyAttributes {
		if value != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// digest the keys and values in a stable order
	digest := sha512.New()
	for _, key := range keys {
		digest.Write([]byte(key))
		digest.Write([]byte("\r"))

		digest.Write([]byte(keyAttributes[key]))
		digest.Write([]byte("\n"))
	}

	sum := make([]byte, digest.Size())
	digest.Sum(sum[:0])
	return fmt.Sprintf("%02x", sum)
}

package vaulted_test

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

var (
	region        = "au-mars-42"
	somethingElse = "something-else"

	testVault = vaulted.Vault{
		AWSKey: &vaulted.AWSKey{
			AWSCredentials: vaulted.AWSCredentials{
				ID:     "id",
				Secret: "secret",
				Region: &region,
			},

			MFA:                     "mfa",
			Role:                    "role",
			ForgoTempCredGeneration: false,
		},

		Vars: map[string]string{
			"TESTING": "testing",
		},

		SSHKeys: map[string]string{
			"SSH_TESTING": "ssh_testing",
		},
	}
)

type UniqSessionKeys map[string]bool

func (u UniqSessionKeys) IsUniq(vault *vaulted.Vault) bool {
	sessionKey := vaulted.VaultSessionCacheKey(vault)
	if _, exists := u[sessionKey]; exists {
		return false
	}

	u[sessionKey] = true
	return true
}

func TestVaultSessionCacheKey(t *testing.T) {
	u := make(UniqSessionKeys)

	var vault vaulted.Vault
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for empty vault")
	}

	vault = testVault
	vault.AWSKey.ID = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key ID")
	}

	vault = testVault
	vault.AWSKey.Secret = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key secret")
	}

	vault = testVault
	vault.AWSKey.Region = &somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key region")
	}

	vault = testVault
	vault.AWSKey.MFA = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key MFA")
	}

	vault = testVault
	vault.AWSKey.Role = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key role")
	}

	vault = testVault
	vault.AWSKey.ForgoTempCredGeneration = true
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered AWS key ID")
	}

	vault = testVault
	vault.Vars["TESTING"] = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered var")
	}

	vault = testVault
	vault.SSHKeys["SSH_TESTING"] = somethingElse
	if !u.IsUniq(&vault) {
		t.Error("Failed to generate unique key for altered SSH key")
	}
}

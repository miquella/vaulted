package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestDump(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{
		AWSKey: &vaulted.AWSKey{
			ID:     "id",
			Secret: "secret",
			MFA:    "mfa",
			Role:   "role",
		},
		Vars: map[string]string{
			"VAR1": "TESTING",
			"VAR2": "ANOTHER",
		},
		SSHKeys: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		},
	}

	output := CaptureStdout(func() {
		d := Dump{
			VaultName: "one",
		}
		err := d.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})

	var v vaulted.Vault
	err := json.Unmarshal(output, &v)
	if err != nil {
		t.Fatalf("Failed to read vault: %v", err)
	}

	if !reflect.DeepEqual(steward.Vaults["one"].AWSKey, v.AWSKey) {
		t.Fatalf("Expected: %#v, got: %#v", steward.Vaults["one"].AWSKey, v.AWSKey)
	}

	if !reflect.DeepEqual(steward.Vaults["one"].Vars, v.Vars) {
		t.Fatalf("Expected: %#v, got: %#v", steward.Vaults["one"].Vars, v.Vars)
	}

	if !reflect.DeepEqual(steward.Vaults["one"].SSHKeys, v.SSHKeys) {
		t.Fatalf("Expected: %#v, got: %#v", steward.Vaults["one"].SSHKeys, v.SSHKeys)
	}

	if !reflect.DeepEqual(steward.Vaults["one"].Duration, v.Duration) {
		t.Fatalf("Expected: %#v, got: %#v", steward.Vaults["one"].Duration, v.Duration)
	}
}

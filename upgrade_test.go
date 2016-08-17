package main

import (
	"reflect"
	"testing"

	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
)

func TestUpgrade(t *testing.T) {
	one := cloneVault(&vaulted.Vault{
		Vars: map[string]string{
			"NEW": "ONE",
		},
	})
	two := cloneVault(&vaulted.Vault{
		Vars: map[string]string{
			"OTHER": "LEGACY",
		},
	})

	steward := NewTestSteward()
	steward.Vaults["one"] = cloneVault(one)
	steward.Passwords["one"] = "hablam wookie"
	steward.LegacyPassword = "robby bobby"
	steward.LegacyEnvironments = map[string]legacy.Environment{
		"one": legacy.Environment{
			Vars: map[string]string{
				"OLD": "LEGACY",
			},
		},
		"two": legacy.Environment{
			Vars: map[string]string{
				"OTHER": "LEGACY",
			},
		},
	}

	CaptureStdout(func() {
		u := Upgrade{}
		err := u.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})

	if !reflect.DeepEqual(one, steward.Vaults["one"]) {
		t.Fatalf("Expected: %#v, got %#v", one, steward.Vaults["one"])
	}
	if steward.Passwords["one"] != "hablam wookie" {
		t.Fatalf("Password should not have changed. Expected %v, got %v", "hablam wookie", steward.Passwords["one"])
	}

	if !reflect.DeepEqual(two, steward.Vaults["two"]) {
		t.Fatalf("Expected: %#v, got %#v", two, steward.Vaults["two"])
	}
	if steward.LegacyPassword != steward.Passwords["two"] {
		t.Fatalf("Password not kept. Expected %v, got %v", steward.LegacyPassword, steward.Passwords["two"])
	}
}

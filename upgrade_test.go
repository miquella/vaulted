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

	store := NewTestStore()
	store.Vaults["one"] = cloneVault(one)
	store.Passwords["one"] = "hablam wookie"
	store.LegacyPassword = "robby bobby"
	store.LegacyEnvironments = map[string]legacy.Environment{
		"one": {
			Vars: map[string]string{
				"OLD": "LEGACY",
			},
		},
		"two": {
			Vars: map[string]string{
				"OTHER": "LEGACY",
			},
		},
	}

	CaptureStdout(func() {
		u := Upgrade{}
		err := u.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})

	if !reflect.DeepEqual(one, store.Vaults["one"]) {
		t.Fatalf("Expected: %#v, got %#v", one, store.Vaults["one"])
	}
	if store.Passwords["one"] != "hablam wookie" {
		t.Fatalf("Password should not have changed. Expected %v, got %v", "hablam wookie", store.Passwords["one"])
	}

	if !reflect.DeepEqual(two, store.Vaults["two"]) {
		t.Fatalf("Expected: %#v, got %#v", two, store.Vaults["two"])
	}
	if store.LegacyPassword != store.Passwords["two"] {
		t.Fatalf("Password not kept. Expected %v, got %v", store.LegacyPassword, store.Passwords["two"])
	}
}

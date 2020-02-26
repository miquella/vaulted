package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestLoad(t *testing.T) {
	v := &vaulted.Vault{
		AWSKey: &vaulted.AWSKey{
			AWSCredentials: vaulted.AWSCredentials{
				ID:     "id",
				Secret: "secret",
			},
			MFA:  "mfa",
			Role: "role",
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
	b, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	store := NewTestStore()
	WriteStdin(b, func() {
		l := Load{
			VaultName: "one",
		}
		err = l.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})

	if !store.VaultExists("one") {
		t.Fatal("The 'one' vault does not exist")
	}

	if !reflect.DeepEqual(v.AWSKey, store.Vaults["one"].AWSKey) {
		t.Fatalf("Expected: %#v, got: %#v", v.AWSKey, store.Vaults["one"].AWSKey)
	}

	if !reflect.DeepEqual(v.Vars, store.Vaults["one"].Vars) {
		t.Fatalf("Expected: %#v, got: %#v", v.Vars, store.Vaults["one"].Vars)
	}

	if !reflect.DeepEqual(v.SSHKeys, store.Vaults["one"].SSHKeys) {
		t.Fatalf("Expected: %#v, got: %#v", v.SSHKeys, store.Vaults["one"].SSHKeys)
	}

	if !reflect.DeepEqual(v.Duration, store.Vaults["one"].Duration) {
		t.Fatalf("Expected: %#v, got: %#v", v.Duration, store.Vaults["one"].Duration)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestLoad(t *testing.T) {
	// Save/restore stdin
	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()

	v := &vaulted.Vault{
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
	b, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	// Write to stdin
	r, w, _ := os.Pipe()
	defer r.Close()
	os.Stdin = r
	go func() {
		vr := bytes.NewReader(b)
		io.Copy(w, vr)
		w.Close()
	}()

	steward := NewTestSteward()
	l := Load{
		VaultName: "one",
	}
	err = l.Run(steward)
	if err != nil {
		t.Fatal(err)
	}

	if !steward.VaultExists("one") {
		t.Fatal("The 'one' vault does not exist")
	}

	if !reflect.DeepEqual(v.AWSKey, steward.Vaults["one"].AWSKey) {
		t.Fatalf("Expected: %#v, got: %#v", v.AWSKey, steward.Vaults["one"].AWSKey)
	}

	if !reflect.DeepEqual(v.Vars, steward.Vaults["one"].Vars) {
		t.Fatalf("Expected: %#v, got: %#v", v.Vars, steward.Vaults["one"].Vars)
	}

	if !reflect.DeepEqual(v.SSHKeys, steward.Vaults["one"].SSHKeys) {
		t.Fatalf("Expected: %#v, got: %#v", v.SSHKeys, steward.Vaults["one"].SSHKeys)
	}

	if !reflect.DeepEqual(v.Duration, steward.Vaults["one"].Duration) {
		t.Fatalf("Expected: %#v, got: %#v", v.Duration, steward.Vaults["one"].Duration)
	}
}

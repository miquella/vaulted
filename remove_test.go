package main

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestRemove(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{}
	steward.Vaults["two"] = &vaulted.Vault{}

	CaptureStdout(func() {
		r := Remove{
			VaultNames: []string{"one"},
		}
		err := r.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
	if steward.VaultExists("one") {
		t.Fatal("The vault 'one' was not removed")
	}

	CaptureStdout(func() {
		r := Remove{
			VaultNames: []string{"one", "two", "three"},
		}
		err := r.Run(steward)
		if err == nil {
			t.Fatal("Expected an error removing 'one' and 'three', but was successful instead")
		}
		exiterr, ok := err.(ErrorWithExitCode)
		if !ok {
			t.Fatalf("Expected ErrorWithExitCode, got %#v", err)
		}
		if exiterr.ExitCode != 2 {
			t.Fatalf("Expected ExitCode: 2, got ExitCode: %v", exiterr.ExitCode)
		}
	})
	if steward.VaultExists("two") {
		t.Fatal("Still expected 'two' to be removed")
	}
}

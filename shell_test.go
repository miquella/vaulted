package main

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestShell(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{}

	WriteStdin([]byte{}, func() {
		s := Shell{
			VaultName: "one",
		}
		err := s.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
}

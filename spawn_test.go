package main

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestSpawn(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{}

	CaptureStdout(func() {
		s := Spawn{
			VaultName: "one",
			Command:   []string{"go", "version"},
		}
		err := s.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
}

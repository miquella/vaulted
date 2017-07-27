package main

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestSpawn(t *testing.T) {
	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{}

	CaptureStdout(func() {
		s := Spawn{
			VaultName: "one",
			Command:   []string{"go", "version"},
		}
		err := s.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})
}

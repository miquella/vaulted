package main

import (
	"strings"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestSpawn(t *testing.T) {
	spawnRefreshVar := `GOPATH="/vaulted"`

	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{}
	store.Vaults["one"].Vars = map[string]string{
		"GOPATH": "/vaulted",
	}

	CaptureStdout(func() {
		s := Spawn{
			SessionOptions: SessionOptions{
				VaultName: "one",
			},
			Command: []string{"go", "version"},
		}
		err := s.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})

	output := CaptureStdout(func() {
		s := Spawn{
			SessionOptions: SessionOptions{
				VaultName: "one",
				Refresh:   true,
			},
			Command:       []string{"go", "env"},
			DisplayStatus: true,
		}
		err := s.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(string(output), spawnRefreshVar) {
		t.Errorf("Incorrect output!\nExpected to contain:\n\"%s\"\ngot:\n\"%s\"", spawnRefreshVar, output)
	}
}

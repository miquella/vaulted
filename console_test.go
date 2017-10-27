package main

import (
	"testing"
	"time"

	"github.com/miquella/vaulted/lib"
)

func TestConsole(t *testing.T) {
	var c Console = Console{
		VaultName: "one",
	}
	var err error
	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{
		AWSKey: &vaulted.AWSKey{},
	}
	store.Vaults["one"].Duration = ConsoleMinDuration
	err = c.Run(store)
	if err != ErrNoCredentialsFound {
		t.Error("No credentials provided, should have caused an ErrNoCredentialsFound")
	}

	store.Vaults["one"].AWSKey.AWSCredentials = vaulted.AWSCredentials{
		ID:     "id",
		Secret: "secret",
	}
	store.Vaults["one"].Duration = 10 * time.Minute
	err = c.Run(store)
	if err != ErrInvalidDuration {
		t.Error("Invalid vault duration, should have caused an ErrInvalidDuration")
	}

	store.Vaults["one"].AWSKey.AWSCredentials = vaulted.AWSCredentials{
		ID:     "id",
		Secret: "secret",
		Token:  "token",
	}
	store.Vaults["one"].Duration = ConsoleMinDuration
	err = c.Run(store)
	if err != ErrInvalidTemporaryCredentials {
		t.Error("Temporary session credentials provided, should have caused an invalid temp credentials error")
	}
}

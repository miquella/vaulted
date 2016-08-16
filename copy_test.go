package main

import (
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestCopy(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["old"] = &vaulted.Vault{
		Vars: map[string]string{
			"TEST": "SUCCESSFUL",
		},
	}

	c := Copy{
		OldVaultName: "old",
		NewVaultName: "new",
	}
	err := c.Run(steward)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := steward.Vaults["new"]
	if !ok {
		t.Fatal("The vault was not copied")
	}

	if v.Vars == nil || v.Vars["TEST"] != "SUCCESSFUL" {
		t.Fatal("The vault contents were not copied")
	}

	if steward.Passwords["old"] == steward.Passwords["new"] {
		t.Fatal("Passwords should be different, but aren't!")
	}
}

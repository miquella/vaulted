package main

import (
	"testing"

	"github.com/miquella/vaulted/v3/lib"
)

func TestCopy(t *testing.T) {
	store := NewTestStore()
	store.Vaults["old"] = &vaulted.Vault{
		Vars: map[string]string{
			"TEST": "SUCCESSFUL",
		},
	}
	store.Passwords["old"] = "one old password"

	c := Copy{
		OldVaultName: "old",
		NewVaultName: "new",
	}
	err := c.Run(store)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := store.Vaults["new"]
	if !ok {
		t.Fatal("The vault was not copied")
	}

	if v.Vars == nil || v.Vars["TEST"] != "SUCCESSFUL" {
		t.Fatal("The vault contents were not copied")
	}

	if store.Passwords["old"] == store.Passwords["new"] {
		t.Fatal("Passwords should be different, but aren't!")
	}
}

func TestCopyToSelf(t *testing.T) {
	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{
		Vars: map[string]string{
			"TEST": "SUCCESSFUL",
		},
	}
	store.Passwords["one"] = "one old password"

	c := Copy{
		OldVaultName: "one",
		NewVaultName: "one",
	}
	err := c.Run(store)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := store.Vaults["one"]
	if !ok {
		t.Fatal("The vault was not preserved")
	}

	if v.Vars == nil || v.Vars["TEST"] != "SUCCESSFUL" {
		t.Fatal("The vault contents were not preserved")
	}

	if store.Passwords["one"] == "one old password" {
		t.Fatal("Passwords should be different, but aren't!")
	}
}

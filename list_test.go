package main

import (
	"bytes"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestList(t *testing.T) {
	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{}
	store.Vaults["two"] = &vaulted.Vault{}

	output := CaptureStdout(func() {
		l := List{}
		err := l.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})

	expected := []byte("one\ntwo\n")
	if bytes.Compare(output, expected) != 0 {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestListWithActive(t *testing.T) {
	store := NewTestStore()
	store.Vaults["first"] = &vaulted.Vault{}
	store.Vaults["second"] = &vaulted.Vault{}
	store.Vaults["third"] = &vaulted.Vault{}

	output := CaptureStdout(func() {
		l := List{
			Active: "second",
		}
		err := l.Run(store)
		if err != nil {
			t.Fatal(err)
		}
	})

	expected := []byte("first\nsecond (active)\nthird\n")
	if bytes.Compare(output, expected) != 0 {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

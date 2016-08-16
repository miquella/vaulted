package main

import (
	"bytes"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestList(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{}
	steward.Vaults["two"] = &vaulted.Vault{}

	output := CaptureStdout(func() {
		l := List{}
		err := l.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})

	expected := []byte("one\ntwo\n")
	if bytes.Compare(output, expected) != 0 {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

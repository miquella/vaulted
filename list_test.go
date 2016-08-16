package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/miquella/vaulted/lib"
)

func TestList(t *testing.T) {
	// Save/restore stdout
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()

	// Capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	captured := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		captured <- buf.Bytes()
		close(captured)
	}()

	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{}
	steward.Vaults["two"] = &vaulted.Vault{}

	l := List{}
	err := l.Run(steward)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	output := <-captured
	expected := []byte("one\ntwo\n")
	if bytes.Compare(output, expected) != 0 {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

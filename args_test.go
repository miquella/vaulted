package main

import (
	"testing"
)

func TestParseArgs_Copy(t *testing.T) {
	good := [][]string{
		[]string{"cp", "one", "two"},
		[]string{"copy", "One", "Two"},
	}

	for _, args := range good {
		cmd, err := ParseArgs(args)
		if err != nil {
			t.Fatalf("Expected %v to parse, error: %v", args, err)
		}

		if cmd == nil {
			t.Fatalf("Expected %v to parse", args)
		}

		if _, ok := cmd.(*Copy); !ok {
			t.Fatalf("Expected %v to produce a Copy command", args)
		}
	}

	bad := [][]string{
		[]string{"cp", "one"},
		[]string{"cp", "one", "two", "three"},
		[]string{"copy", "one"},
		[]string{"copy", "one", "two", "three"},
	}

	for _, args := range bad {
		_, err := ParseArgs(args)
		if err == nil {
			t.Fatalf("Expected %v to fail to parse", args)
		}
	}
}

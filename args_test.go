package main

import (
	"reflect"
	"testing"
)

type parseCase struct {
	Args    []string
	Command Command
}

var (
	goodParseCases = []parseCase{
		// Copy
		{
			Args: []string{"cp", "one", "two"},
			Command: &Copy{
				OldVaultName: "one",
				NewVaultName: "two",
			},
		},
		{
			Args: []string{"copy", "one", "two"},
			Command: &Copy{
				OldVaultName: "one",
				NewVaultName: "two",
			},
		},

		// Dump
		{
			Args: []string{"dump", "one"},
			Command: &Dump{
				VaultName: "one",
			},
		},

		// List
		{
			Args:    []string{"ls"},
			Command: &List{},
		},
		{
			Args:    []string{"list"},
			Command: &List{},
		},

		// Load
		{
			Args: []string{"load", "one"},
			Command: &Load{
				VaultName: "one",
			},
		},

		// Remove
		{
			Args: []string{"rm", "one"},
			Command: &Remove{
				VaultNames: []string{"one"},
			},
		},
		{
			Args: []string{"rm", "one", "two", "three", "four"},
			Command: &Remove{
				VaultNames: []string{"one", "two", "three", "four"},
			},
		},

		// Shell
		{
			Args: []string{"shell", "one"},
			Command: &Shell{
				VaultName: "one",
			},
		},
	}

	badParseCases = []parseCase{
		// Copy
		{
			Args: []string{"cp"},
		},
		{
			Args: []string{"cp", "one"},
		},
		{
			Args: []string{"cp", "one", "two", "three"},
		},
		{
			Args: []string{"copy"},
		},
		{
			Args: []string{"copy", "one"},
		},
		{
			Args: []string{"copy", "one", "two", "three"},
		},

		// Dump
		{
			Args: []string{"dump"},
		},
		{
			Args: []string{"dump", "one", "two"},
		},

		// List
		{
			Args: []string{"ls", "one"},
		},
		{
			Args: []string{"list", "one"},
		},

		// Load
		{
			Args: []string{"load"},
		},
		{
			Args: []string{"load", "one", "two"},
		},

		// Remove
		{
			Args: []string{"rm"},
		},

		// Shell
		{
			Args: []string{"shell"},
		},
		{
			Args: []string{"shell", "one", "two"},
		},
	}
)

type parseExpectation struct {
	Args    []string
	Command Command
}

func TestParseArgs(t *testing.T) {
	for _, good := range goodParseCases {
		cmd, err := ParseArgs(good.Args)
		if err != nil {
			t.Fatalf("Failed to parse '%v': %v", good.Args, err)
		}

		if !reflect.DeepEqual(good.Command, cmd) {
			t.Fatalf("Expected command: %#v, got: %#v", good.Command, cmd)
		}
	}

	for _, bad := range badParseCases {
		_, err := ParseArgs(bad.Args)
		if err == nil {
			t.Fatalf("Expected '%v' to fail to parse", bad.Args)
		}
	}
}

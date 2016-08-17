package main

import (
	"os"
	"reflect"
	"testing"
)

type parseCase struct {
	Args    []string
	Command Command
}

var (
	goodParseCases = []parseCase{
		// Spawn
		{
			Args: []string{"-n", "one"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"/bin/fish", "--login"},
			},
		},
		{
			Args: []string{"-n", "one", "-i"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"/bin/fish", "--login"},
			},
		},
		{
			Args: []string{"-i", "-n", "one"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"/bin/fish", "--login"},
			},
		},
		{
			Args: []string{"-in", "one"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"/bin/fish", "--login"},
			},
		},
		{
			Args: []string{"-n", "one", "some", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"some", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "some", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"some", "command"},
			},
		},

		// Add
		{
			Args: []string{"add", "one"},
			Command: &Edit{
				VaultName: "one",
			},
		},

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

		// Edit
		{
			Args: []string{"edit", "one"},
			Command: &Edit{
				VaultName: "one",
			},
		},

		// Env
		{
			Args: []string{"env", "one"},
			Command: &Env{
				VaultName: "one",
				Shell:     "fish",
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
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"/bin/fish", "--login"},
			},
		},

		// Upgrade
		{
			Args:    []string{"upgrade"},
			Command: &Upgrade{},
		},
	}

	badParseCases = []parseCase{
		// Spawn
		{
			Args: []string{"-i"},
		},
		{
			Args: []string{"-n", "one", "-i", "some", "command"},
		},
		{
			Args: []string{"-n", "one", "-i", "--", "some", "command"},
		},
		{
			Args: []string{"-n", "one", "some", "--", "command"},
		},

		// Add
		{
			Args: []string{"add"},
		},
		{
			Args: []string{"add", "one", "two"},
		},

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

		// Edit
		{
			Args: []string{"edit"},
		},
		{
			Args: []string{"edit", "one", "two"},
		},

		// Env
		{
			Args: []string{"env"},
		},
		{
			Args: []string{"env", "one", "two"},
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

		// Upgrade
		{
			Args: []string{"upgrade", "one"},
		},

		// Misc
		{
			Args: []string{"bobby"},
		},
	}
)

type parseExpectation struct {
	Args    []string
	Command Command
}

func TestParseArgs(t *testing.T) {
	shell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", shell)
	os.Setenv("SHELL", "/bin/fish")

	for _, good := range goodParseCases {
		var cmd Command
		var err error
		CaptureStdout(func() {
			cmd, err = ParseArgs(good.Args)
		})
		if err != nil {
			t.Errorf("Failed to parse '%v': %v", good.Args, err)
		}

		if !reflect.DeepEqual(good.Command, cmd) {
			t.Errorf("Expected command: %#v, got: %#v", good.Command, cmd)
		}
	}

	for _, bad := range badParseCases {
		_, err := ParseArgs(bad.Args)
		if err == nil {
			t.Errorf("Expected '%v' to fail to parse", bad.Args)
		}
	}
}

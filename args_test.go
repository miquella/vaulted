package main

import (
	"os"
	"reflect"
	"testing"
)

type parseCase struct {
	Args        []string
	OsArgs      []string
	Environment map[string]string

	Command Command
}

var (
	goodParseCases = []parseCase{
		// Spawn
		{
			Args: []string{"-n", "one"},
			Command: &Spawn{
				VaultName:     "one",
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-n", "one", "-i"},
			Command: &Spawn{
				VaultName:     "one",
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-i", "-n", "one"},
			Command: &Spawn{
				VaultName:     "one",
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-in", "one"},
			Command: &Spawn{
				VaultName:     "one",
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
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
			Args: []string{"-n", "one", "command", "--command-flag"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"command", "--command-flag"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "some", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"some", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "some", "--", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"some", "--", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "some", "--", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"some", "--", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "--", "some", "--", "command"},
			Command: &Spawn{
				VaultName: "one",
				Command:   []string{"--", "some", "--", "command"},
			},
		},

		// Add
		{
			Args: []string{"add", "one"},
			Command: &Edit{
				VaultName: "one",
			},
		},
		{
			Args:    []string{"add", "--help"},
			Command: &Help{Subcommand: "add"},
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
		{
			Args:    []string{"copy", "--help"},
			Command: &Help{Subcommand: "copy"},
		},
		{
			Args:    []string{"cp", "--help"},
			Command: &Help{Subcommand: "cp"},
		},

		// Dump
		{
			Args: []string{"dump", "one"},
			Command: &Dump{
				VaultName: "one",
			},
		},
		{
			Args:    []string{"dump", "--help"},
			Command: &Help{Subcommand: "dump"},
		},

		// Edit
		{
			Args: []string{"edit", "one"},
			Command: &Edit{
				VaultName: "one",
			},
		},
		{
			Args:    []string{"edit", "--help"},
			Command: &Help{Subcommand: "edit"},
		},

		// Env
		{
			Args:   []string{"env", "one"},
			OsArgs: []string{"vaulted", "env", "one"},
			Command: &Env{
				VaultName:     "one",
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env one",
			},
		},
		{
			Args:    []string{"env", "--help"},
			Command: &Help{Subcommand: "env"},
		},
		{
			Args:   []string{"env", "foo", "--format", "json"},
			OsArgs: []string{"vaulted", "env", "foo", "--format", "json"},
			Command: &Env{
				VaultName:     "foo",
				DetectedShell: "fish",
				Format:        "json",
				Command:       "vaulted env foo --format json",
			},
		},

		// Help
		{
			Args:    []string{"help", "add"},
			Command: &Help{Subcommand: "add"},
		},
		{
			Args:    []string{"help", "cp"},
			Command: &Help{Subcommand: "cp"},
		},
		{
			Args:    []string{"help", "copy"},
			Command: &Help{Subcommand: "copy"},
		},
		{
			Args:    []string{"help", "dump"},
			Command: &Help{Subcommand: "dump"},
		},
		{
			Args:    []string{"help", "edit"},
			Command: &Help{Subcommand: "edit"},
		},
		{
			Args:    []string{"help", "env"},
			Command: &Help{Subcommand: "env"},
		},
		{
			Args:    []string{"help", "list"},
			Command: &Help{Subcommand: "list"},
		},
		{
			Args:    []string{"help", "ls"},
			Command: &Help{Subcommand: "ls"},
		},
		{
			Args:    []string{"help", "load"},
			Command: &Help{Subcommand: "load"},
		},
		{
			Args:    []string{"help", "rm"},
			Command: &Help{Subcommand: "rm"},
		},
		{
			Args:    []string{"help", "shell"},
			Command: &Help{Subcommand: "shell"},
		},
		{
			Args:    []string{"help", "upgrade"},
			Command: &Help{Subcommand: "upgrade"},
		},
		{
			Args:    []string{"-h"},
			Command: &Help{},
		},
		{
			Args:    []string{"--help"},
			Command: &Help{},
		},

		// List
		{
			Args:    []string{"ls"},
			Command: &List{},
		},
		{
			Args: []string{"ls"},
			Environment: map[string]string{
				"VAULTED_ENV": "active-env",
			},
			Command: &List{
				Active: "active-env",
			},
		},
		{
			Args:    []string{"list"},
			Command: &List{},
		},
		{
			Args: []string{"list"},
			Environment: map[string]string{
				"VAULTED_ENV": "active-env",
			},
			Command: &List{
				Active: "active-env",
			},
		},
		{
			Args:    []string{"list", "--help"},
			Command: &Help{Subcommand: "list"},
		},
		{
			Args:    []string{"ls", "--help"},
			Command: &Help{Subcommand: "ls"},
		},

		// Load
		{
			Args: []string{"load", "one"},
			Command: &Load{
				VaultName: "one",
			},
		},
		{
			Args:    []string{"load", "--help"},
			Command: &Help{Subcommand: "load"},
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
		{
			Args:    []string{"rm", "--help"},
			Command: &Help{Subcommand: "rm"},
		},

		// Shell
		{
			Args: []string{"shell", "one"},
			Command: &Spawn{
				VaultName:     "one",
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args:    []string{"shell", "--help"},
			Command: &Help{Subcommand: "shell"},
		},

		// Upgrade
		{
			Args:    []string{"upgrade"},
			Command: &Upgrade{},
		},
		{
			Args:    []string{"upgrade", "--help"},
			Command: &Help{Subcommand: "upgrade"},
		},

		//Version
		{
			Args:    []string{"version"},
			Command: &Version{},
		},
		{
			Args:    []string{"-V"},
			Command: &Version{},
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
			Args: []string{},
		},
		{
			Args: []string{"bobby"},
		},
		{
			Args: []string{"blah", "--help"},
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
		// Temporarily set environment variables
		savedEnv := make(map[string]string)
		for key, value := range good.Environment {
			savedEnv[key] = os.Getenv(key)
			os.Setenv(key, value)
		}

		if good.OsArgs != nil {
			origArgs := os.Args
			os.Args = good.OsArgs
			defer func() { os.Args = origArgs }()
		}

		var cmd Command
		var err error
		CaptureStdout(func() {
			cmd, err = ParseArgs(good.Args)
		})
		if err != nil {
			t.Errorf("Failed to parse %#v: %v", good.Args, err)
		}

		if !reflect.DeepEqual(good.Command, cmd) {
			t.Errorf("Expected command: %#v, got: %#v", good.Command, cmd)
		}

		// Restore environment variables
		for key, value := range savedEnv {
			os.Setenv(key, value)
		}
	}

	for _, bad := range badParseCases {
		// Temporarily set environment variables
		savedEnv := make(map[string]string)
		for key, value := range bad.Environment {
			savedEnv[key] = os.Getenv(key)
			os.Setenv(key, value)
		}

		_, err := ParseArgs(bad.Args)
		if err == nil {
			t.Errorf("Expected %#v to fail to parse", bad.Args)
		}

		// Restore environment variables
		for key, value := range savedEnv {
			os.Setenv(key, value)
		}
	}
}

package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/miquella/vaulted/edit"
)

type parseCase struct {
	Args        []string
	OsArgs      []string
	Environment map[string]string

	Command Command
}

var (
	twoHourDuration = time.Hour * 2
	goodParseCases  = []parseCase{
		// Spawn
		{
			Args: []string{"-n", "one"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-n", "one", "-i"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-i", "-n", "one"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-in", "one"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"-n", "one", "some", "command"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"some", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "command", "--command-flag"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"command", "--command-flag"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "some", "command"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"some", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "some", "--", "command"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"some", "--", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "some", "--", "command"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"some", "--", "command"},
			},
		},
		{
			Args: []string{"-n", "one", "--", "--", "some", "--", "command"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"--", "some", "--", "command"},
			},
		},

		// Add
		{
			Args: []string{"add", "one"},
			Command: &edit.Edit{
				New:       true,
				VaultName: "one",
			},
		},
		{
			Args: []string{"new", "one"},
			Command: &edit.Edit{
				New:       true,
				VaultName: "one",
			},
		},
		{
			Args: []string{"create", "one"},
			Command: &edit.Edit{
				New:       true,
				VaultName: "one",
			},
		},
		{
			Args:    []string{"add", "--help"},
			Command: &Help{Subcommand: "add"},
		},
		{
			Args:    []string{"new", "--help"},
			Command: &Help{Subcommand: "new"},
		},
		{
			Args:    []string{"create", "--help"},
			Command: &Help{Subcommand: "create"},
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
			Command: &edit.Edit{
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
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env one",
			},
		},
		{
			Args:   []string{"env", "--assume", "arn:something:or:other", "one"},
			OsArgs: []string{"vaulted", "env", "--assume", "arn:something:or:other", "one"},
			Command: &Env{
				SessionOptions: SessionOptions{
					VaultName: "one",
					Role:      "arn:something:or:other",
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env --assume arn:something:or:other one",
			},
		},
		{
			Args:   []string{"env", "--assume", "arn:something:or:other"},
			OsArgs: []string{"vaulted", "env", "--assume", "arn:something:or:other"},
			Command: &Env{
				SessionOptions: SessionOptions{
					Role: "arn:something:or:other",
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env --assume arn:something:or:other",
			},
		},
		{
			Args:   []string{"env", "foo", "--format", "json"},
			OsArgs: []string{"vaulted", "env", "foo", "--format", "json"},
			Command: &Env{
				SessionOptions: SessionOptions{
					VaultName: "foo",
				},
				DetectedShell: "fish",
				Format:        "json",
				Command:       "vaulted env foo --format json",
			},
		},
		{
			Args:   []string{"env", "foo", "--no-session"},
			OsArgs: []string{"vaulted", "env", "foo", "--no-session"},
			Command: &Env{
				SessionOptions: SessionOptions{
					VaultName: "foo",
					NoSession: true,
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env foo --no-session",
			},
		},
		{
			Args:   []string{"env", "foo", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			OsArgs: []string{"vaulted", "env", "foo", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			Command: &Env{
				SessionOptions: SessionOptions{
					VaultName:    "foo",
					Role:         "arn:something:or:other",
					RoleDuration: &twoHourDuration,
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env foo --assume arn:something:or:other --assume-duration 2h",
			},
		},
		{
			Args:   []string{"env", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			OsArgs: []string{"vaulted", "env", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			Command: &Env{
				SessionOptions: SessionOptions{
					Role:         "arn:something:or:other",
					RoleDuration: &twoHourDuration,
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env --assume arn:something:or:other --assume-duration 2h",
			},
		},
		{
			Args:   []string{"env", "foo", "--assume-duration", "2h"},
			OsArgs: []string{"vaulted", "env", "foo", "--assume-duration", "2h"},
			Command: &Env{
				SessionOptions: SessionOptions{
					VaultName:    "foo",
					RoleDuration: &twoHourDuration,
				},
				DetectedShell: "fish",
				Format:        "shell",
				Command:       "vaulted env foo --assume-duration 2h",
			},
		},
		{
			Args:    []string{"env", "--help"},
			Command: &Help{Subcommand: "env"},
		},

		// Exec
		{
			Args:   []string{"exec", "one", "cmd", "cmd2"},
			OsArgs: []string{"vaulted", "exec", "one", "cmd", "cmd2"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"cmd", "cmd2"},
			},
		},
		{
			Args:   []string{"exec", "one", "--", "cmd", "cmd2"},
			OsArgs: []string{"vaulted", "exec", "one", "--", "cmd", "cmd2"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command: []string{"cmd", "cmd2"},
			},
		},
		{
			Args:   []string{"exec", "--assume", "arn:some:thing", "one", "cmd", "cmd2"},
			OsArgs: []string{"vaulted", "exec", "--assume", "arn:some:thing", "one", "cmd", "cmd2"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
					Role:      "arn:some:thing",
				},
				Command: []string{"cmd", "cmd2"},
			},
		},
		{
			Args:   []string{"exec", "--assume", "arn:some:thing", "--", "cmd", "cmd2"},
			OsArgs: []string{"vaulted", "exec", "--assume", "arn:some:thing", "cmd", "cmd2"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					Role: "arn:some:thing",
				},
				Command: []string{"cmd", "cmd2"},
			},
		},
		{
			Args:   []string{"exec", "--no-session", "one", "cmd", "cmd2"},
			OsArgs: []string{"vaulted", "exec", "--no-session", "one", "cmd", "cmd2"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
					NoSession: true,
				},
				Command: []string{"cmd", "cmd2"},
			},
		},
		{
			Args:    []string{"exec", "--help"},
			Command: &Help{Subcommand: "exec"},
		},

		// Help
		{
			Args:    []string{"help", "add"},
			Command: &Help{Subcommand: "add"},
		},
		{
			Args:    []string{"help", "new"},
			Command: &Help{Subcommand: "new"},
		},
		{
			Args:    []string{"help", "create"},
			Command: &Help{Subcommand: "create"},
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
			Args:    []string{"help", "exec"},
			Command: &Help{Subcommand: "exec"},
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
			Args:    []string{"help", "passwd"},
			Command: &Help{Subcommand: "passwd"},
		},
		{
			Args:    []string{"help", "password"},
			Command: &Help{Subcommand: "password"},
		},
		{
			Args:    []string{"help", "rm"},
			Command: &Help{Subcommand: "rm"},
		},
		{
			Args:    []string{"help", "remove"},
			Command: &Help{Subcommand: "remove"},
		},
		{
			Args:    []string{"help", "delete"},
			Command: &Help{Subcommand: "delete"},
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

		// Passwd
		{
			Args: []string{"passwd", "one"},
			Command: &Copy{
				OldVaultName: "one",
				NewVaultName: "one",
			},
		},
		{
			Args: []string{"password", "one"},
			Command: &Copy{
				OldVaultName: "one",
				NewVaultName: "one",
			},
		},
		{
			Args:    []string{"passwd", "--help"},
			Command: &Help{Subcommand: "passwd"},
		},
		{
			Args:    []string{"password", "--help"},
			Command: &Help{Subcommand: "password"},
		},

		// Remove
		{
			Args: []string{"rm", "one"},
			Command: &Remove{
				VaultNames: []string{"one"},
			},
		},
		{
			Args: []string{"remove", "one"},
			Command: &Remove{
				VaultNames: []string{"one"},
			},
		},
		{
			Args: []string{"delete", "one"},
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
			Args: []string{"remove", "one", "two", "three", "four"},
			Command: &Remove{
				VaultNames: []string{"one", "two", "three", "four"},
			},
		},
		{
			Args: []string{"delete", "one", "two", "three", "four"},
			Command: &Remove{
				VaultNames: []string{"one", "two", "three", "four"},
			},
		},
		{
			Args:    []string{"rm", "--help"},
			Command: &Help{Subcommand: "rm"},
		},
		{
			Args:    []string{"remove", "--help"},
			Command: &Help{Subcommand: "remove"},
		},
		{
			Args:    []string{"delete", "--help"},
			Command: &Help{Subcommand: "delete"},
		},

		// Shell
		{
			Args: []string{"shell", "one"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "--assume", "arn:something:or:other"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					Role: "arn:something:or:other",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "--assume", "arn:something:or:other", "one"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "one",
					Role:      "arn:something:or:other",
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "foo", "--no-session"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName: "foo",
					NoSession: true,
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "foo", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName:    "foo",
					Role:         "arn:something:or:other",
					RoleDuration: &twoHourDuration,
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "--assume", "arn:something:or:other", "--assume-duration", "2h"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					Role:         "arn:something:or:other",
					RoleDuration: &twoHourDuration,
				},
				Command:       []string{"/bin/fish", "--login"},
				DisplayStatus: true,
			},
		},
		{
			Args: []string{"shell", "foo", "--assume-duration", "2h"},
			Command: &Spawn{
				SessionOptions: SessionOptions{
					VaultName:    "foo",
					RoleDuration: &twoHourDuration,
				},
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
		{
			Args: []string{"env", "--no-session"},
		},
		{
			Args: []string{"env", "one", "--no-session", "--assume", "arn:blah:blah"},
		},
		{
			Args: []string{"env", "one", "--no-session", "--refresh"},
		},
		{
			Args: []string{"env", "--assume-duration", "2h"},
		},

		// Exec
		{
			// no arguments provided
			Args: []string{"exec"},
		},
		{
			// no command provided
			Args: []string{"exec", "one"},
		},
		{
			// no command provided
			Args: []string{"exec", "--assume", "arn:some:thing"},
		},
		{
			// must provide vault name or --assume
			Args: []string{"exec", "--", "cmd"},
		},
		{
			// may not provide -- without args following
			Args: []string{"exec", "one", "--"},
		},
		{
			// may not provide cmd args before and after the --
			Args: []string{"exec", "one", "cmd", "--", "cmd2"},
		},
		{
			// may not provide --no-session without a vault name
			Args: []string{"exec", "--no-session", "--", "cmd"},
		},
		{
			// may not provide both --no-session and --assume
			Args: []string{"exec", "one", "--no-session", "--assume", "arn:some:thing", "cmd"},
		},
		{
			// may not provide both --no-session and --refresh
			Args: []string{"exec", "one", "--no-session", "--refresh", "cmd"},
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

		// Passwd
		{
			Args: []string{"passwd"},
		},
		{
			Args: []string{"password"},
		},
		{
			Args: []string{"passwd", "one", "two"},
		},
		{
			Args: []string{"password", "one", "two"},
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
		{
			Args: []string{"shell", "--no-session"},
		},
		{
			Args: []string{"shell", "one", "--no-session", "--assume", "arn:blah:blah"},
		},
		{
			Args: []string{"shell", "one", "--no-session", "--refresh"},
		},
		{
			Args: []string{"shell", "--assume-duration", "2h"},
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
	// backup & nuke env vars we don't want
	savedEnv := make(map[string]string)
	for _, env := range os.Environ() {
		envPieces := strings.SplitN(env, "=", 2)
		if envPieces[0] == "SHELL" || strings.HasPrefix(envPieces[0], "VAULTED") {
			savedEnv[envPieces[0]] = envPieces[1]
			os.Unsetenv(envPieces[0])
		}
	}
	defer func() {
		for key, value := range savedEnv {
			os.Setenv(key, value)
		}
	}()

	// set SHELL to our control case
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

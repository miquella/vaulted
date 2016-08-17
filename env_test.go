package main

import (
	"os"
	"testing"

	"github.com/miquella/vaulted/lib"
)

var (
	envFishOutput = `set -x ONE "111111";
set -x THREE "333";
set -x TWO "222";
`
	envFishOutputWithHint = `# To load these variables into your shell, execute:
#   eval (vaulted env one)
set -x ONE "111111";
set -x THREE "333";
set -x TWO "222";
`

	envShOutput = `export ONE="111111"
export THREE="333"
export TWO="222"
`
	envShOutputWithHint = `# To load these variables into your shell, execute:
#   eval $(vaulted env one)
export ONE="111111"
export THREE="333"
export TWO="222"
`
)

func TestEng(t *testing.T) {
	steward := NewTestSteward()
	steward.Vaults["one"] = &vaulted.Vault{
		Vars: map[string]string{
			"TWO":   "222",
			"ONE":   "111111",
			"THREE": "333",
		},
	}

	output := CaptureStdout(func() {
		e := Env{
			VaultName: "one",
			Shell:     "fish",
		}
		err := e.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
	if string(output) != envFishOutput {
		t.Fatalf("Incorrect output: %s", output)
	}

	output = CaptureStdout(func() {
		args := os.Args
		os.Args = []string{"vaulted", "env", "one"}
		defer func() { os.Args = args }()

		e := Env{
			VaultName: "one",
			Shell:     "fish",
			UsageHint: true,
		}
		err := e.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
	if string(output) != envFishOutputWithHint {
		t.Fatalf("Incorrect output: %s", output)
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName: "one",
			Shell:     "sh",
		}
		err := e.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
	if string(output) != envShOutput {
		t.Fatalf("Incorrect output: %s", output)
	}

	output = CaptureStdout(func() {
		args := os.Args
		os.Args = []string{"vaulted", "env", "one"}
		defer func() { os.Args = args }()

		e := Env{
			VaultName: "one",
			Shell:     "sh",
			UsageHint: true,
		}
		err := e.Run(steward)
		if err != nil {
			t.Fatal(err)
		}
	})
	if string(output) != envShOutputWithHint {
		t.Fatalf("Incorrect output: %s", output)
	}
}

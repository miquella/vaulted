package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/miquella/vaulted/lib"
)

var (
	envFishOutput = `# To load these variables into your shell, execute:
#   vaulted env one | source
set -gx ONE "111111";
set -gx THREE "333";
set -gx TWO "222";
set -gx VAULTED_ENV "one";
set -gx VAULTED_ENV_EXPIRATION "2006-01-02T22:04:05Z";
`
	envFishOutputWithPermCreds = `# To load these variables into your shell, execute:
#   vaulted env one | source
set -e AWS_SECURITY_TOKEN;
set -e AWS_SESSION_TOKEN;
set -gx AWS_ACCESS_KEY_ID "aws-key-id";
set -gx AWS_SECRET_ACCESS_KEY "aws-secret-key";
set -gx ONE "111111";
set -gx THREE "333";
set -gx TWO "222";
set -gx VAULTED_ENV "one";
set -gx VAULTED_ENV_EXPIRATION "2006-01-02T22:04:05Z";
`
	envFishOutputNonInteractive = `set -gx ONE "111111";
set -gx THREE "333";
set -gx TWO "222";
set -gx VAULTED_ENV "one";
set -gx VAULTED_ENV_EXPIRATION "2006-01-02T22:04:05Z";
`
	envFishOutputWithPermCredsNonInteractive = `set -e AWS_SECURITY_TOKEN;
set -e AWS_SESSION_TOKEN;
set -gx AWS_ACCESS_KEY_ID "aws-key-id";
set -gx AWS_SECRET_ACCESS_KEY "aws-secret-key";
set -gx ONE "111111";
set -gx THREE "333";
set -gx TWO "222";
set -gx VAULTED_ENV "one";
set -gx VAULTED_ENV_EXPIRATION "2006-01-02T22:04:05Z";
`

	envShOutput = `# To load these variables into your shell, execute:
#   eval "$(vaulted env one)"
export ONE="111111"
export THREE="333"
export TWO="222"
export VAULTED_ENV="one"
export VAULTED_ENV_EXPIRATION="2006-01-02T22:04:05Z"
`
	envShOutputWithPermCreds = `# To load these variables into your shell, execute:
#   eval "$(vaulted env one)"
unset AWS_SECURITY_TOKEN
unset AWS_SESSION_TOKEN
export AWS_ACCESS_KEY_ID="aws-key-id"
export AWS_SECRET_ACCESS_KEY="aws-secret-key"
export ONE="111111"
export THREE="333"
export TWO="222"
export VAULTED_ENV="one"
export VAULTED_ENV_EXPIRATION="2006-01-02T22:04:05Z"
`

	envShOutputNonInteractive = `export ONE="111111"
export THREE="333"
export TWO="222"
export VAULTED_ENV="one"
export VAULTED_ENV_EXPIRATION="2006-01-02T22:04:05Z"
`
	envShOutputWithPermCredsNonInteractive = `unset AWS_SECURITY_TOKEN
unset AWS_SESSION_TOKEN
export AWS_ACCESS_KEY_ID="aws-key-id"
export AWS_SECRET_ACCESS_KEY="aws-secret-key"
export ONE="111111"
export THREE="333"
export TWO="222"
export VAULTED_ENV="one"
export VAULTED_ENV_EXPIRATION="2006-01-02T22:04:05Z"
`

	envJSONOutput = `{
  "ONE": "111111",
  "THREE": "333",
  "TWO": "222",
  "VAULTED_ENV": "one",
  "VAULTED_ENV_EXPIRATION": "2006-01-02T22:04:05Z"
}
`
	envRefreshOutput = `export ONE="111111"
export THREE="333"
export TWO="222"
export VAULTED_ENV="one"
export VAULTED_ENV_EXPIRATION="2006-01-02T22:04:06Z"
`

	envCustom = "[AWS_SECURITY_TOKEN AWS_SESSION_TOKEN]"
)

func TestEnv(t *testing.T) {
	store := NewTestStore()
	store.Vaults["one"] = &vaulted.Vault{
		Vars: map[string]string{
			"TWO":   "222",
			"ONE":   "111111",
			"THREE": "333",
		},
	}

	output := CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "fish",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   true,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envFishOutput {
		t.Error(failureMessage(envFishOutput, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   true,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envShOutput {
		t.Error(failureMessage(envShOutput, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "fish",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   false,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envFishOutputNonInteractive {
		t.Error(failureMessage(envFishOutputNonInteractive, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   false,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envShOutputNonInteractive {
		t.Error(failureMessage(envShOutputNonInteractive, output))
	}

	output = CaptureStdout(func() {
		args := os.Args
		os.Args = []string{"vaulted", "env", "one", "--format", "json"}
		defer func() { os.Args = args }()

		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "json",
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envJSONOutput {
		t.Error(failureMessage(envJSONOutput, output))
	}

	// cached session
	store.Sessions["one"] = &vaulted.Session{
		SessionVersion: vaulted.SessionVersion,

		Name:       "one",
		Expiration: time.Unix(1136239445, 0),

		AWSCreds: &vaulted.AWSCredentials{
			ID:     "aws-key-id",
			Secret: "aws-secret-key",
		},
		Vars: map[string]string{
			"TWO":   "222",
			"ONE":   "111111",
			"THREE": "333",
		},
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "fish",
			Format:        "fish",
			Command:       "vaulted env one",
			Interactive:   true,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envFishOutputWithPermCreds {
		t.Error(failureMessage(envFishOutputWithPermCreds, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   true,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "fish",
			Format:        "fish",
			Command:       "vaulted env one",
			Interactive:   false,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envFishOutputWithPermCredsNonInteractive {
		t.Error(failureMessage(envFishOutputWithPermCredsNonInteractive, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "shell",
			Command:       "vaulted env one",
			Interactive:   false,
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})

	if string(output) != envShOutputWithPermCredsNonInteractive {
		t.Error(failureMessage(envShOutputWithPermCredsNonInteractive, output))
	}

	output = CaptureStdout(func() {
		args := os.Args
		os.Args = []string{"vaulted", "env", "one"}
		defer func() { os.Args = args }()

		e := Env{
			VaultName: "one",
			Format:    "{{ .Unset }}",
		}
		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})

	if string(output) != envCustom {
		t.Error(failureMessage(envCustom, output))
	}

	output = CaptureStdout(func() {
		e := Env{
			VaultName:     "one",
			DetectedShell: "sh",
			Format:        "shell",
			Command:       "vaulted env one --refresh",
			Refresh:       true,
			Interactive:   false,
		}

		err := e.Run(store)
		if err != nil {
			t.Error(err)
		}
	})
	if string(output) != envRefreshOutput {
		t.Error(failureMessage(envRefreshOutput, output))
	}

}

func failureMessage(expected string, got []byte) string {
	return fmt.Sprintf("Incorrect output!\nExpected:\n\"%s\"\ngot:\n\"%s\"", expected, got)
}

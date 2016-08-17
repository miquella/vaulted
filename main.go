package main

import (
	"errors"
	"fmt"
	"os"
)

type ErrorWithExitCode struct {
	error
	ExitCode int
}

var (
	ErrNoError = errors.New("")
)

func main() {
	command, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(255)
	}

	if command == nil {
		// no command, display usage help instead
		PrintUsage()
		os.Exit(255)
	}

	steward := &TTYSteward{}
	err = command.Run(steward)
	if err != nil {
		exiterr, ok := err.(ErrorWithExitCode)
		if !ok || exiterr.error != ErrNoError {
			fmt.Fprintln(os.Stderr, err)
		}
		if ok {
			os.Exit(exiterr.ExitCode)
		} else {
			os.Exit(1)
		}
	}
}

func PrintUsage() {
	fmt.Fprintln(os.Stderr, "NAME")
	fmt.Fprintln(os.Stderr, "    vaulted - spawn environments from securely stored secrets")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "SYNOPSIS")
	fmt.Fprintln(os.Stderr, "    vaulted -n VAULT [-i]")
	fmt.Fprintln(os.Stderr, "    vaulted -n VAULT [--] CMD")
	fmt.Fprintln(os.Stderr, "    vaulted COMMAND [args...]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "DESCRIPTION")
	fmt.Fprintln(os.Stderr, "    If no *COMMAND* is provided, `vaulted` either spawns *CMD* (if provided) or spawns an interactive shell.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "COMMANDS")
	fmt.Fprintln(os.Stderr, "    add        - Interactively creates the content of a new vault.")
	fmt.Fprintln(os.Stderr, "    cp / copy  - Copies the content of a vault and saves it as a new vault with a new password.")
	fmt.Fprintln(os.Stderr, "    dump       - Writes the content of a vault to stdout as JSON.")
	fmt.Fprintln(os.Stderr, "    edit       - Interactively edits the content of an existing vault.")
	fmt.Fprintln(os.Stderr, "    env        - Outputs shell commands that load secrets for a vault into the shell.")
	fmt.Fprintln(os.Stderr, "    load       - Uses JSON provided to stdin to create or replace the content of a vault.")
	fmt.Fprintln(os.Stderr, "    ls / list  - Lists all vaults.")
	fmt.Fprintln(os.Stderr, "    rm         - Removes existing vaults.")
	fmt.Fprintln(os.Stderr, "    shell      - Starts an interactive shell with the secrets for the vault loaded into the shell.")
	fmt.Fprintln(os.Stderr, "    upgrade    - Upgrades legacy vaults to the current vault format.")

}

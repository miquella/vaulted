package main

import (
	"fmt"
	"os"

	"github.com/miquella/vaulted/lib"
)

func main() {
	// omit the command name that is passed to VaultedCLI
	var cli VaultedCLI
	if len(os.Args) > 0 {
		cli = VaultedCLI(os.Args[1:])
	}

	cli.Run()
}

type VaultedCLI []string

func (cli VaultedCLI) Run() {
	if len(cli) == 0 {
		os.Exit(255)
	}

	switch cli[0] {
	case "list", "ls":
		cli.List()

	case "rm":
		cli.Remove()

	default:
		os.Exit(255)
	}
}

func (cli VaultedCLI) List() {
	vaults, err := vaulted.ListVaults()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to list vaults: %v", err))
		os.Exit(1)
	}

	for _, vault := range vaults {
		fmt.Fprintln(os.Stdout, vault)
	}
}

func (cli VaultedCLI) Remove() {
	if len(cli) <= 1 {
		fmt.Fprintln(os.Stderr, "You must specify which vaults to remove")
		os.Exit(255)
	}

	failures := 0
	for _, name := range cli[1:] {
		err := vaulted.RemoveVault(name)
		if err != nil {
			failures++
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %v", name, err))
		}
	}

	os.Exit(failures)
}

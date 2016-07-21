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
		os.Exit(1)
	}

	switch cli[0] {
	case "list", "ls":
		cli.List()

	default:
		os.Exit(1)
	}
}

func (cli VaultedCLI) List() {
	vaults, err := vaulted.ListVaults()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to list vaults: %v", err))
		os.Exit(2)
	}

	for _, vault := range vaults {
		fmt.Fprintln(os.Stdout, vault)
	}
}

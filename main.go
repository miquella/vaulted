package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/bgentry/speakeasy"
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

func getPassword() string {
	password, err := speakeasy.Ask("Password: ")
	if err != nil {
		os.Exit(1)
	}
	return password
}

type VaultedCLI []string

func (cli VaultedCLI) Run() {
	if len(cli) == 0 {
		os.Exit(255)
	}

	switch cli[0] {
	case "cat":
		cli.Cat()

	case "list", "ls":
		cli.List()

	case "rm":
		cli.Remove()

	case "shell":
		cli.Shell()

	default:
		os.Exit(255)
	}
}

func (cli VaultedCLI) Cat() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a single vault to cat")
		os.Exit(255)
	}

	password := getPassword()
	vault, err := vaulted.OpenVault(password, cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	staticVars, err := vault.GetEnvVars(nil, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var keys []string
	for key, _ := range staticVars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s=%s", key, staticVars[key]))
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

func (cli VaultedCLI) Shell() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a single vault to spawn a shell with")
		os.Exit(255)
	}

	password := getPassword()
	vault, err := vaulted.OpenVault(password, cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code, err := vault.Spawn([]string{os.Getenv("SHELL"), "--login"}, map[string]string{"VAULTED_ENV": cli[1]})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(*code)
}

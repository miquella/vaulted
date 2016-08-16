package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"github.com/miquella/vaulted/lib/legacy"
	"github.com/spf13/pflag"
)

var (
	ErrUnknownShell = errors.New("Unknown shell")
)

func main() {
	command, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(255)
	}

	if command != nil {
		steward := &TTYSteward{}
		err := command.Run(steward)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// omit the command name that is passed to VaultedCLI
	var cli VaultedCLI
	if len(os.Args) > 0 {
		cli = VaultedCLI(os.Args[1:])
	}

	cli.Run()
}

func openEnvironment(name string) (password string, env *vaulted.Environment, err error) {
	password = os.Getenv("VAULTED_PASSWORD")
	if password != "" {
		env, err = vaulted.GetEnvironment(name, password)
	} else {
		for i := 0; i < 3; i++ {
			password, err = ask.HiddenAsk("Password: ")
			if err != nil {
				break
			}

			env, err = vaulted.GetEnvironment(name, password)
			if err != vaulted.ErrInvalidPassword {
				break
			}
		}
	}
	return
}

func openVault(name string) (password string, vault *vaulted.Vault, err error) {
	password = os.Getenv("VAULTED_PASSWORD")
	if password != "" {
		vault, err = vaulted.OpenVault(password, name)
	} else {
		for i := 0; i < 3; i++ {
			password, err = ask.HiddenAsk("Password: ")
			if err != nil {
				break
			}

			vault, err = vaulted.OpenVault(password, name)
			if err != vaulted.ErrInvalidPassword {
				break
			}
		}
	}
	return
}

func openLegacyVault() (password string, environments map[string]legacy.Environment, err error) {
	legacyVault, err := legacy.ReadVault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	password = os.Getenv("VAULTED_PASSWORD")
	if password != "" {
		environments, err = legacyVault.DecryptEnvironments(password)
	} else {
		for i := 0; i < 3; i++ {
			password, err = ask.HiddenAsk("Legacy Password: ")
			if err != nil {
				break
			}

			environments, err = legacyVault.DecryptEnvironments(password)
			if err != legacy.ErrInvalidPassword {
				break
			}
		}
	}
	return
}

type VaultedCLI []string

func (cli VaultedCLI) Run() {
	if len(cli) == 0 {
		cli.PrintUsage()
		os.Exit(255)
	}

	switch cli[0] {
	case "add", "edit":
		cli.Edit()

	case "dump":
		cli.Dump()

	case "env":
		cli.Env()

	case "list", "ls":
		cli.List()

	case "load":
		cli.Load()

	case "rm":
		cli.Remove()

	case "shell":
		cli.Shell()

	case "upgrade":
		cli.Upgrade()

	case "help":
		cli.PrintUsage()
		os.Exit(255)

	default:
		if strings.HasPrefix(cli[0], "-") {
			cli.Spawn()
		} else {
			fmt.Fprintf(os.Stderr, "Invalid command: %s\n", cli[0])
			cli.PrintUsage()
			os.Exit(255)
		}
	}
}

func (cli VaultedCLI) PrintUsage() {
	fmt.Fprintln(os.Stderr, "USAGE:")
	fmt.Fprintln(os.Stderr, "  vaulted -n VAULT [--] CMD    - Spawn CMD in the VAULT environment")
	fmt.Fprintln(os.Stderr, "  vaulted -n VAULT [-i]        - Spawn an interactive shell in the VAULT environment")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  vaulted ls                   - List all vaults")
	fmt.Fprintln(os.Stderr, "  vaulted add VAULT            - Interactively add the VAULT")
	fmt.Fprintln(os.Stderr, "  vaulted edit VAULT           - Interactively edit the VAULT")
	fmt.Fprintln(os.Stderr, "  vaulted cp VAULT NEWVAULT    - Creates a copy of VAULT as NEWVAULT")
	fmt.Fprintln(os.Stderr, "  vaulted rm VAULT [VAULT...]  - Remove the VAULT environment(s)")
	fmt.Fprintln(os.Stderr, "  vaulted env VAULT            - Displays the environment variables for the VAULT environment")
	fmt.Fprintln(os.Stderr, "  vaulted shell VAULT          - Spawn an interactive shell in the VAULT environment")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  vaulted dump VAULT           - Dump the VAULT in JSON format")
	fmt.Fprintln(os.Stderr, "  vaulted load VAULT           - Load the VAULT from JSON format")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  vaulted upgrade              - Upgrade from a legacy vaulted format")
}

func (cli VaultedCLI) Dump() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a vault to dump")
		os.Exit(255)
	}

	_, vault, err := openVault(cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jvault, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for len(jvault) > 0 {
		n, err := os.Stdout.Write(jvault)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		jvault = jvault[n:]
	}
}

func (cli VaultedCLI) Env() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a vault for which to get the environment")
		os.Exit(255)
	}

	_, env, err := openEnvironment(cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// detect the correct shell
	shell, err := detectShell()
	if err == ErrUnknownShell {
		shell = "sh"
	}

	usageHint := ""
	setVar := ""
	quoteReplacement := "\""
	switch shell {
	case "fish":
		usageHint = "# To load these variables into your shell, execute:\n#   eval (%s)"
		setVar = "set -x %s \"%s\";"
		quoteReplacement = "\\\""
	default:
		usageHint = "# To load these variables into your shell, execute:\n#   eval $(%s)"
		setVar = "export %s=\"%s\""
		quoteReplacement = "\\\""
	}

	// sort the vars
	var keys []string
	for key, _ := range env.Vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// display the vars using the format string for the shell
	displayUsageHint := true
	fi, err := os.Stdout.Stat()
	if err == nil {
		if fi.Mode()&os.ModeCharDevice == 0 {
			displayUsageHint = false
		}
	}
	if displayUsageHint {
		fmt.Fprintln(os.Stdout, fmt.Sprintf(usageHint, strings.Join(os.Args, " ")))
	}

	for _, key := range keys {
		fmt.Fprintln(os.Stdout, fmt.Sprintf(setVar, key, strings.Replace(env.Vars[key], "\"", quoteReplacement, -1)))
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

func (cli VaultedCLI) Load() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a vault to load")
		os.Exit(255)
	}

	jvault, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	vault := &vaulted.Vault{}
	err = json.Unmarshal(jvault, vault)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	password, err := ask.HiddenAsk("New Password: ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = vaulted.SealVault(password, cli[1], vault)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "You must specify a vault to spawn a shell with")
		os.Exit(255)
	}

	currentVaultedEnv := os.Getenv("VAULTED_ENV")
	if currentVaultedEnv != "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Refusing to spawn a new shell when already in environment '%s'.", currentVaultedEnv))
		os.Exit(255)
	}

	_, env, err := openEnvironment(cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code, err := env.Spawn([]string{os.Getenv("SHELL"), "--login"}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(*code)
}

func (cli VaultedCLI) Spawn() {
	spawnFlags := pflag.NewFlagSet("spawn", pflag.ContinueOnError)
	spawnFlags.SetInterspersed(false)

	name := spawnFlags.StringP("name", "n", "", "Name of the vault to spawn")
	interactive := spawnFlags.BoolP("interactive", "i", false, "Spawn an interactive shell")
	force := spawnFlags.BoolP("force", "f", false, "Bypass protective checks and force spawning of the environment")
	help := spawnFlags.Bool("help", false, "Show usage help")
	err := spawnFlags.Parse([]string(cli))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(255)
	}

	if *help {
		cli.PrintUsage()
		os.Exit(255)
	}

	currentVaultedEnv := os.Getenv("VAULTED_ENV")
	if !*force && currentVaultedEnv != "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Refusing to spawn a new environment when already in environment '%s'. Use --force to override.", currentVaultedEnv))
		os.Exit(255)
	}

	if spawnFlags.ArgsLenAtDash() > 0 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unknown argument(s): %v", spawnFlags.Args()[:spawnFlags.ArgsLenAtDash()]))
		os.Exit(255)
	}

	if *name == "" {
		*name = os.Getenv("VAULTED_DEFAULT_ENV")
	}

	if *name == "" {
		fmt.Fprintln(os.Stderr, "A vault must be specified when spawning")
		os.Exit(255)
	}

	_, env, err := openEnvironment(cli[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cmd []string
	if *interactive || len(spawnFlags.Args()) == 0 {
		cmd = append(cmd, os.Getenv("SHELL"), "--login")
	}
	cmd = append(cmd, spawnFlags.Args()...)

	code, err := env.Spawn(cmd, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(*code)
}

func (cli VaultedCLI) Upgrade() {
	password, environments, err := openLegacyVault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// collect the current list of vaults (so we don't overwrite any)
	vaults, _ := vaulted.ListVaults()
	existingVaults := map[string]bool{}
	for _, name := range vaults {
		existingVaults[name] = true
	}

	failed := 0
	for name, env := range environments {
		if existingVaults[name] {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: skipped (vault already exists)", name))
			continue
		}

		vault := vaulted.Vault{
			Vars: env.Vars,
		}
		err = vaulted.SealVault(password, name, &vault)
		if err != nil {
			failed++
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %v", name, err))
		} else {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: upgraded", name))
		}
	}

	os.Exit(failed)
}

func detectShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell), nil
	}

	return "", ErrUnknownShell
}

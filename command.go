package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/miquella/vaulted/lib"
	"github.com/spf13/pflag"
)

var (
	ErrSubcommandRequired          = ErrorWithExitCode{errors.New("A subcommand must be specified. See 'vaulted --help' for details."), EX_USAGE_ERROR}
	ErrTooManyArguments            = ErrorWithExitCode{errors.New("too many arguments provided"), EX_USAGE_ERROR}
	ErrNotEnoughArguments          = ErrorWithExitCode{errors.New("not enough arguments provided"), EX_USAGE_ERROR}
	ErrVaultNameRequired           = ErrorWithExitCode{errors.New("A vault name must be specified"), EX_USAGE_ERROR}
	ErrMixingCommandAndInteractive = ErrorWithExitCode{errors.New("Cannot mix an interactive shell with command arguments"), EX_USAGE_ERROR}

	ErrUnknownShell = errors.New("Unknown shell")
)

type Command interface {
	Run(store vaulted.Store) error
}

func ParseArgs(args []string) (Command, error) {
	command, err := parseArgs(args)
	if err == pflag.ErrHelp {
		if HelpAliases[args[0]] == "" {
			return parseHelpArgs(nil)
		} else {
			return parseHelpArgs(args)
		}
	}

	// If arguments fail to parse for any reason, it's a usage error
	if err != nil {
		if _, ok := err.(ErrorWithExitCode); !ok {
			err = ErrorWithExitCode{err, EX_USAGE_ERROR}
		}
	}

	return command, err
}

func parseArgs(args []string) (Command, error) {
	flag := spawnFlagSet()
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.Changed("version") {
		return &Version{}, nil
	}

	if flag.Changed("name") || flag.Changed("interactive") {
		return parseSpawnArgs(args)
	}

	// Parse command
	commandArgs := flag.Args()
	if len(commandArgs) == 0 || flag.ArgsLenAtDash() == 0 {
		return nil, ErrSubcommandRequired
	}

	if flag.ArgsLenAtDash() > -1 {
		commandArgsWithDash := append([]string{}, commandArgs[:flag.ArgsLenAtDash()]...)
		commandArgsWithDash = append(commandArgsWithDash, "--")
		commandArgsWithDash = append(commandArgsWithDash, commandArgs[flag.ArgsLenAtDash():]...)
		commandArgs = commandArgsWithDash
	}

	switch commandArgs[0] {
	case "add", "create", "new":
		return parseAddArgs(commandArgs[1:])

	case "cp", "copy":
		return parseCopyArgs(commandArgs[1:])

	case "dump":
		return parseDumpArgs(commandArgs[1:])

	case "edit":
		return parseEditArgs(commandArgs[1:])

	case "env":
		return parseEnvArgs(commandArgs[1:])

	case "help":
		return parseHelpArgs(commandArgs[1:])

	case "ls", "list":
		return parseListArgs(commandArgs[1:])

	case "load":
		return parseLoadArgs(commandArgs[1:])

	case "rm", "delete", "remove":
		return parseRemoveArgs(commandArgs[1:])

	case "shell":
		return parseShellArgs(commandArgs[1:])

	case "upgrade":
		return parseUpgradeArgs(commandArgs[1:])

	case "version":
		return &Version{}, nil

	default:
		return nil, fmt.Errorf("Unknown command: %s", commandArgs[0])
	}
}

func spawnFlagSet() *pflag.FlagSet {
	flag := pflag.NewFlagSet("vaulted", pflag.ContinueOnError)
	flag.Usage = func() {}
	flag.SetInterspersed(false)
	flag.StringP("name", "n", "", "Name of the vault to use")
	flag.BoolP("interactive", "i", false, "Spawn interactive shell (if -n is used, but no additional arguments a provided, interactive is the default)")
	flag.BoolP("version", "V", false, "Specify current version of Vaulted")
	return flag
}

func parseSpawnArgs(args []string) (Command, error) {
	flag := spawnFlagSet()
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	name, _ := flag.GetString("name")
	interactive, _ := flag.GetBool("interactive")

	if name == "" {
		return nil, ErrVaultNameRequired
	}

	if flag.ArgsLenAtDash() > 0 {
		return nil, fmt.Errorf("Unknown arguments: %s", strings.Join(flag.Args()[:flag.ArgsLenAtDash()], " "))
	}

	if interactive && flag.NArg() > 0 {
		return nil, ErrMixingCommandAndInteractive
	}

	currentVaultedEnv := os.Getenv("VAULTED_ENV")
	if currentVaultedEnv != "" {
		return nil, fmt.Errorf("Refusing to spawn a new shell when already in environment '%s'.", currentVaultedEnv)
	}

	s := &Spawn{}
	s.VaultName = name
	if interactive || flag.NArg() == 0 {
		s.Command = interactiveShellCommand()
		s.DisplayStatus = true
	} else {
		s.Command = flag.Args()
	}
	return s, nil
}

func parseAddArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("add", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	e := &Edit{}
	e.New = true
	e.VaultName = flag.Arg(0)
	return e, nil
}

func parseCopyArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("copy", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 2 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 2 {
		return nil, ErrTooManyArguments
	}

	c := &Copy{}
	c.OldVaultName = flag.Arg(0)
	c.NewVaultName = flag.Arg(1)
	return c, nil
}

func parseDumpArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("dump", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	d := &Dump{}
	d.VaultName = flag.Arg(0)
	return d, nil
}

func parseEditArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("edit", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	e := &Edit{}
	e.VaultName = flag.Arg(0)
	return e, nil
}

func parseEnvArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("env", pflag.ContinueOnError)
	flag.String("format", "shell", "Specify what built in format to output variables in (shell, sh, fish, json) or a text template. Default: shell")
	flag.String("assume", "", "Role to assume")
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	vaultName := ""
	assume, _ := flag.GetString("assume")

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	if flag.NArg() == 1 {
		vaultName = flag.Arg(0)
	} else if flag.NArg() < 1 && assume == "" {
		return nil, ErrNotEnoughArguments
	}

	shell, err := detectShell()
	if err == ErrUnknownShell {
		shell = "sh"
	}

	e := &Env{}
	e.VaultName = vaultName
	e.Role = assume
	e.DetectedShell = shell
	e.Command = strings.Join(os.Args, " ")

	e.Interactive = true
	fi, err := os.Stdout.Stat()
	if err == nil {
		if fi.Mode()&os.ModeCharDevice == 0 {
			e.Interactive = false
		}
	}

	e.Format, _ = flag.GetString("format")
	return e, nil
}

func parseHelpArgs(args []string) (Command, error) {
	h := Help{}
	if len(args) > 0 {
		h.Subcommand = args[0]
	}
	return &h, nil
}

func parseListArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("list", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	return &List{Active: os.Getenv("VAULTED_ENV")}, nil
}

func parseLoadArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("load", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	l := &Load{}
	l.VaultName = flag.Arg(0)
	return l, nil
}

func parseRemoveArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("remove", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	r := &Remove{}
	r.VaultNames = flag.Args()
	return r, nil
}

func parseShellArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("shell", pflag.ContinueOnError)
	flag.String("assume", "", "Role to assume")
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	vaultName := ""
	assume, _ := flag.GetString("assume")

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	if flag.NArg() == 1 {
		vaultName = flag.Arg(0)

		currentVaultedEnv := os.Getenv("VAULTED_ENV")
		if currentVaultedEnv != "" {
			return nil, fmt.Errorf("Refusing to spawn a new shell when already in environment '%s'.", currentVaultedEnv)
		}
	} else if flag.NArg() < 1 && assume == "" {
		return nil, ErrNotEnoughArguments
	}

	s := &Spawn{}
	s.VaultName = vaultName
	s.Role = assume
	s.Command = interactiveShellCommand()
	s.DisplayStatus = true
	return s, nil
}

func parseUpgradeArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("upgrade", pflag.ContinueOnError)
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	return &Upgrade{}, nil
}

func interactiveShellCommand() []string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	return []string{shell, "--login"}
}

func detectShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell), nil
	}

	return "", ErrUnknownShell
}

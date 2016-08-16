package main

import (
	"errors"

	"github.com/spf13/pflag"
)

var (
	ErrTooManyArguments   = errors.New("too many arguments provided")
	ErrNotEnoughArguments = errors.New("not enough arguments provided")
)

func ParseArgs(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, nil
	}

	switch args[0] {
	case "cp", "copy":
		return parseCopyArgs(args[1:])

	case "dump":
		return parseDumpArgs(args[1:])

	case "ls", "list":
		return parseListArgs(args[1:])

	case "load":
		return parseLoadArgs(args[1:])

	case "rm":
		return parseRemoveArgs(args[1:])

	default:
		return nil, nil
	}
}

func parseCopyArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("copy", pflag.ContinueOnError)
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

func parseListArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("list", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	return &List{}, nil
}

func parseLoadArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("load", pflag.ContinueOnError)
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

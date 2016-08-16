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

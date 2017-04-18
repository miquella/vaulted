package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	EX_USAGE_ERROR     = 64
	EX_DATA_ERROR      = 65
	EX_TEMPORARY_ERROR = 79
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
	if err == nil {
		steward := &TTYSteward{}
		err = command.Run(steward)
	}

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

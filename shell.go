package main

import (
	"os"
)

type Shell struct {
	VaultName string
}

func (s *Shell) Run(steward Steward) error {
	_, env, err := steward.GetEnvironment(s.VaultName, nil)
	if err != nil {
		return err
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	code, err := env.Spawn([]string{shell, "--login"}, nil)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

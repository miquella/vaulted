package main

import (
	"errors"
	"fmt"
)

type Remove struct {
	VaultNames []string
}

func (r *Remove) Run(steward Steward) error {
	failures := 0
	for _, name := range r.VaultNames {
		err := steward.RemoveVault(name)
		if err != nil {
			failures++
			fmt.Printf("%s: %v\n", name, err)
		}
	}

	if failures > 0 {
		return ErrorWithExitCode{
			errors.New("Vault could not be removed"),
			failures,
		}
	}

	return nil
}

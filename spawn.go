package main

import (
	"math"
	"strings"

	"github.com/miquella/vaulted/lib"
)

type Spawn struct {
	VaultNames []string
	Command    []string
}

func (s *Spawn) Run(steward Steward) error {
	var envs []*vaulted.Environment
	for _, name := range s.VaultNames {
		_, env, err := steward.GetEnvironment(name, nil)
		if err != nil {
			return err
		}

		envs = append(envs, env)
	}

	mergedEnv := vaulted.Environment{
		Expiration: math.MaxInt64,
		Vars:       map[string]string{},
		SSHKeys:    map[string]string{},
	}
	for _, env := range envs {
		if env.Expiration < mergedEnv.Expiration {
			mergedEnv.Expiration = env.Expiration
		}

		for key, value := range env.Vars {
			mergedEnv.Vars[key] = value
		}

		for key, value := range env.SSHKeys {
			mergedEnv.SSHKeys[key] = value
		}
	}

	mergedEnv.Vars["VAULTED_ENV"] = strings.Join(s.VaultNames, "\n")

	code, err := mergedEnv.Spawn(s.Command, nil)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

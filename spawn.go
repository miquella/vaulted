package main

import (
	"fmt"
	"time"

	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
)

type Spawn struct {
	VaultName string
	Role      string

	Command       []string
	DisplayStatus bool
}

func (s *Spawn) Run(steward Steward) error {
	env, err := s.getEnvironment(steward)
	if err != nil {
		return err
	}

	timeRemaining := env.Expiration.Sub(time.Now())
	timeRemaining = time.Second * time.Duration(timeRemaining.Seconds())
	if s.DisplayStatus {
		ask.Print(fmt.Sprintf("%s — expires: %s (%s remaining)\n", s.VaultName, env.Expiration.Format("2 Jan 2006 15:04 MST"), timeRemaining))
	}

	code, err := env.Spawn(s.Command)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

func (s *Spawn) getEnvironment(steward Steward) (*vaulted.Environment, error) {
	var err error

	// default environment
	env := &vaulted.Environment{
		Expiration: time.Now().Add(time.Hour),
	}

	if s.VaultName != "" {
		// get specific environment
		_, env, err = steward.GetEnvironment(s.VaultName, nil)
		if err != nil {
			return nil, err
		}
	}

	if s.Role != "" {
		return env.Assume(s.Role)
	}

	return env, nil
}

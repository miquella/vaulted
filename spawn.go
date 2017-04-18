package main

import (
	"fmt"
	"time"

	"github.com/miquella/ask"
)

type Spawn struct {
	VaultName     string
	Command       []string
	DisplayStatus bool
}

func (s *Spawn) Run(steward Steward) error {
	_, env, err := steward.GetEnvironment(s.VaultName, nil)
	if err != nil {
		return err
	}

	if s.DisplayStatus {
		expiration := env.Expiration
		timeRemaining := expiration.Sub(time.Now())
		timeRemaining = time.Second * time.Duration(timeRemaining.Seconds())
		ask.Print(fmt.Sprintf("%s — expires: %s (%s remaining)\n", s.VaultName, expiration.Format("2 Jan 2006 15:04 MST"), timeRemaining))
	}

	code, err := env.Spawn(s.Command, nil)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

package main

import (
	"fmt"
	"time"

	"github.com/miquella/ask"
	"github.com/miquella/ssh-proxy-agent/lib/proxyagent"

	"github.com/miquella/vaulted/lib"
)

type Spawn struct {
	SessionOptions

	Command       []string
	DisplayStatus bool
}

func (s *Spawn) Run(store vaulted.Store) error {
	session, err := GetSessionWithOptions(store, &s.SessionOptions)
	if err != nil {
		return err
	}

	timeRemaining := session.Expiration.Sub(time.Now())
	timeRemaining = time.Second * time.Duration(timeRemaining.Seconds())
	if s.DisplayStatus {
		ask.Print(fmt.Sprintf("%s — expires: %s (%s remaining)\n", session.Name, session.Expiration.Format("2 Jan 2006 15:04 MST"), timeRemaining))
	}

	sshAgent, err := proxyagent.SetupAgent(proxyagent.AgentConfig{})
	if err != nil {
		return ErrorWithExitCode{err, 1}
	}

	code, err := session.Spawn(s.Command, sshAgent)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

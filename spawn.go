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
	Refresh       bool
}

func (s *Spawn) Run(store vaulted.Store) error {
	session, err := s.getSession(store)
	if err != nil {
		return err
	}

	timeRemaining := session.Expiration.Sub(time.Now())
	timeRemaining = time.Second * time.Duration(timeRemaining.Seconds())
	if s.DisplayStatus {
		ask.Print(fmt.Sprintf("%s — expires: %s (%s remaining)\n", session.Name, session.Expiration.Format("2 Jan 2006 15:04 MST"), timeRemaining))
	}

	code, err := session.Spawn(s.Command)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}

func (s *Spawn) getSession(store vaulted.Store) (*vaulted.Session, error) {
	var err error

	// default session
	session := DefaultSession()

	if s.VaultName != "" {
		// get specific session
		if s.Refresh {
			session, _, err = store.CreateSession(s.VaultName)
		} else {
			session, _, err = store.GetSession(s.VaultName)
		}
		if err != nil {
			return nil, err
		}
	}

	if s.Role != "" {
		return session.Assume(s.Role)
	}

	return session, nil
}

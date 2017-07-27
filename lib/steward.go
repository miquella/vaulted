package vaulted

import (
	"errors"
)

type Operation int

const (
	OpenOperation Operation = iota
	SealOperation
)

type Steward interface {
	GetMFAToken(name string) (string, error)
	GetPassword(operation Operation, name string) (string, error)
}

type StewardMaxTries interface {
	GetMaxOpenTries() int
}

type StaticSteward struct {
	Password string
	MFAToken *string
}

func NewStaticSteward(password string) *StaticSteward {
	return &StaticSteward{
		Password: password,
	}
}

func NewStaticStewardWithMFA(password, mfaToken string) *StaticSteward {
	return &StaticSteward{
		Password: password,
		MFAToken: &mfaToken,
	}
}

func (s *StaticSteward) GetPassword(operation Operation, name string) (string, error) {
	return s.Password, nil
}

func (s *StaticSteward) GetMFAToken(name string) (string, error) {
	if s.MFAToken == nil {
		return "", errors.New("No MFA token available")
	} else {
		return *s.MFAToken, nil
	}
}

package main

import (
	"fmt"
)

const (
	VERSION = "2.0.unstable"
)

type Version struct{}

func (l *Version) Run(steward Steward) error {

	fmt.Printf("Vaulted v%s\n", VERSION)
	return nil
}

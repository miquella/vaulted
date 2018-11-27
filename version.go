package main

import (
	"fmt"

	"github.com/miquella/vaulted/lib"
)

const (
	VERSION = "2.4.unstable"
)

type Version struct{}

func (l *Version) Run(store vaulted.Store) error {
	fmt.Printf("Vaulted v%s\n", VERSION)
	return nil
}

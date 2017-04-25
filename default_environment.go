package main

import (
	"os"
	"time"

	"github.com/miquella/vaulted/lib"
)

func DefaultEnvironment() *vaulted.Environment {
	return &vaulted.Environment{
		Name:       os.Getenv("VAULTED_ENV"),
		Expiration: time.Now().Add(time.Hour).Truncate(time.Second),
	}
}

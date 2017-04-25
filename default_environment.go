package main

import (
	"os"
	"time"

	"github.com/miquella/vaulted/lib"
)

func DefaultSession() *vaulted.Session {
	return &vaulted.Session{
		Name:       os.Getenv("VAULTED_ENV"),
		Expiration: time.Now().Add(time.Hour).Truncate(time.Second),
	}
}

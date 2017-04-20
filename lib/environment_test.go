package vaulted

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestEnvironmentVariables(t *testing.T) {
	e := Environment{
		Expiration: time.Now(),
		Vars: map[string]string{
			"TEST":         "TESTING",
			"ANOTHER_TEST": "TEST TEST",
		},
	}
	var expectedSet map[string]string = map[string]string{
		"ANOTHER_TEST": "TEST TEST",
		"TEST":         "TESTING",
		"VAULTED_ENV_EXPIRATION": e.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset []string

	vars := e.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

func TestEnvironmentVariablesWithPermCreds(t *testing.T) {
	e := Environment{
		Expiration: time.Now(),
		AWSCreds: &AWSCredentials{
			ID:     "an-id",
			Secret: "the-super-sekrit",
		},
		Vars: map[string]string{
			"TEST":         "TESTING",
			"ANOTHER_TEST": "TEST TEST",
		},
	}
	var expectedSet map[string]string = map[string]string{
		"ANOTHER_TEST":          "TEST TEST",
		"AWS_ACCESS_KEY_ID":     e.AWSCreds.ID,
		"AWS_SECRET_ACCESS_KEY": e.AWSCreds.Secret,
		"TEST":                  "TESTING",
		"VAULTED_ENV_EXPIRATION": e.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset []string = []string{
		"AWS_SECURITY_TOKEN",
		"AWS_SESSION_TOKEN",
	}

	vars := e.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	sort.Strings(vars.Unset)
	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

func TestEnvironmentVariablesWithTempCreds(t *testing.T) {
	e := Environment{
		Expiration: time.Now(),
		AWSCreds: &AWSCredentials{
			ID:     "an-id",
			Secret: "the-super-sekrit",
			Token:  "my-affections",
		},
		Vars: map[string]string{
			"TEST":         "TESTING",
			"ANOTHER_TEST": "TEST TEST",
		},
	}
	var expectedSet map[string]string = map[string]string{
		"ANOTHER_TEST":          "TEST TEST",
		"AWS_ACCESS_KEY_ID":     e.AWSCreds.ID,
		"AWS_SECRET_ACCESS_KEY": e.AWSCreds.Secret,
		"AWS_SECURITY_TOKEN":    e.AWSCreds.Token,
		"AWS_SESSION_TOKEN":     e.AWSCreds.Token,
		"TEST":                  "TESTING",
		"VAULTED_ENV_EXPIRATION": e.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset []string

	vars := e.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

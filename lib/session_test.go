package vaulted

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestSessionVariables(t *testing.T) {
	s := Session{
		Name:       "vault",
		Expiration: time.Now(),

		Vars: map[string]string{
			"TEST":         "TESTING",
			"ANOTHER_TEST": "TEST TEST",
		},
	}
	var expectedSet = map[string]string{
		"ANOTHER_TEST":           "TEST TEST",
		"TEST":                   "TESTING",
		"VAULTED_ENV":            s.Name,
		"VAULTED_ENV_EXPIRATION": s.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset []string

	vars := s.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

func TestSessionVariablesWithPermCreds(t *testing.T) {
	s := Session{
		Name:       "vault",
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
	var expectedSet = map[string]string{
		"ANOTHER_TEST":           "TEST TEST",
		"AWS_ACCESS_KEY_ID":      s.AWSCreds.ID,
		"AWS_SECRET_ACCESS_KEY":  s.AWSCreds.Secret,
		"TEST":                   "TESTING",
		"VAULTED_ENV":            s.Name,
		"VAULTED_ENV_EXPIRATION": s.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset = []string{
		"AWS_SECURITY_TOKEN",
		"AWS_SESSION_TOKEN",
	}

	vars := s.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	sort.Strings(vars.Unset)
	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

func TestSessionVariablesWithTempCreds(t *testing.T) {
	s := Session{
		Name:       "vault",
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
	var expectedSet = map[string]string{
		"ANOTHER_TEST":           "TEST TEST",
		"AWS_ACCESS_KEY_ID":      s.AWSCreds.ID,
		"AWS_SECRET_ACCESS_KEY":  s.AWSCreds.Secret,
		"AWS_SECURITY_TOKEN":     s.AWSCreds.Token,
		"AWS_SESSION_TOKEN":      s.AWSCreds.Token,
		"TEST":                   "TESTING",
		"VAULTED_ENV":            s.Name,
		"VAULTED_ENV_EXPIRATION": s.Expiration.UTC().Format(time.RFC3339),
	}
	var expectedUnset []string

	vars := s.Variables()

	if !reflect.DeepEqual(expectedSet, vars.Set) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedSet, vars.Set)
	}

	if !reflect.DeepEqual(expectedUnset, vars.Unset) {
		t.Errorf("Expected: %#v\nGot: %#v\n", expectedUnset, vars.Unset)
	}
}

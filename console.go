package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/miquella/vaulted/lib"
	"github.com/pkg/browser"
)

var (
	ErrInvalidTemporaryCredentials = errors.New("Temporary session credentials found. Console cannot be opened with temporary credentials.\nIf the AWS Key is set in a vault, the permanent credentials can be used by specifying the vault name.")
	ErrInvalidDuration             = errors.New("Console duration must be between 15m and 12h.")
	ErrNoCredentialsFound          = errors.New("No credentials found. Console cannot be opened.")
)

const (
	CONSOLE_URL = "https://console.aws.amazon.com/console/home"
	SIGNIN_URL  = "https://signin.aws.amazon.com/federation"

	MinConsoleDuration = 15 * time.Minute
	MaxConsoleDuration = 12 * time.Hour
)

type Console struct {
	VaultName string
	Role      string
	Duration  time.Duration
}

func (c *Console) Run(store vaulted.Store) error {
	signinToken, err := c.getSigninToken(store)
	if err != nil {
		return err
	}

	// console signin
	signinUrl, _ := url.Parse(SIGNIN_URL)
	loginQuery := make(url.Values)
	loginQuery.Set("Action", "login")
	loginQuery.Set("SigninToken", signinToken)
	loginQuery.Set("Destination", CONSOLE_URL)

	signinUrl.RawQuery = loginQuery.Encode()
	err = browser.OpenURL(signinUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func (c *Console) getSigninToken(store vaulted.Store) (string, error) {
	// Setup default values (may be overwritten by values from vault)
	duration := 1 * time.Hour
	var awsKey vaulted.AWSKey

	// Override defaults with values from specified vault
	if c.VaultName != "" {
		v, _, err := store.OpenVault(c.VaultName)
		if err != nil {
			return "", err
		}

		duration = v.Duration

		if v.AWSKey.Valid() {
			awsKey = *v.AWSKey
			if c.Role != "" {
				awsKey.Role = c.Role
			}
		}
	}

	// If duration was provided through the command line overwrite with that
	if c.Duration != 0 && (c.Duration < 15*time.Minute || c.Duration > 12*time.Hour) {
		return "", ErrInvalidDuration
	} else if c.Duration > 0 {
		duration = c.Duration
	}

	return c.getSigninTokenFromCreds(store, awsKey, duration)
}

func (c *Console) getSigninTokenFromCreds(store vaulted.Store, awsKey vaulted.AWSKey, duration time.Duration) (string, error) {
	// Get creds from environment if no creds loaded from vault
	creds, err := awsKey.AWSCredentials.WithLocalDefault()
	if err != nil && err != credentials.ErrNoValidProvidersFoundInChain {
		return "", err
	} else if err == credentials.ErrNoValidProvidersFoundInChain || !creds.Valid() {
		return "", ErrNoCredentialsFound
	}

	if creds.ValidSession() {
		return "", ErrInvalidTemporaryCredentials
	}

	duration, err = capDuration(duration)
	if err != nil {
		return "", err
	}

	// assume provided role or get a federation token
	if awsKey.Role != "" {
		if awsKey.MFA != "" {
			tokenCode, tokenErr := store.Steward().GetMFAToken(c.VaultName)
			if tokenErr != nil {
				return "", tokenErr
			}
			creds, err = creds.AssumeRoleWithMFA(awsKey.MFA, tokenCode, awsKey.Role, 15*time.Minute)
		} else {
			creds, err = creds.AssumeRole(awsKey.Role, 15*time.Minute)
		}
		if err != nil {
			return "", err
		}
		return creds.GetSigninToken(&duration)
	} else {
		creds, err = creds.GetFederationToken(duration)
		if err != nil {
			return "", err
		}
		return creds.GetSigninToken(nil)
	}
}

func capDuration(duration time.Duration) (time.Duration, error) {
	if duration < MinConsoleDuration {
		return time.Duration(0), ErrInvalidDuration
	}
	if duration > MaxConsoleDuration {
		duration = MaxConsoleDuration
		fmt.Println("Your vault duration is greater than the max console duration.\nCurrent console session duration set to 12 hours.")
	}
	return duration, nil
}

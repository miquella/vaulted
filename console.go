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
	ConsoleURL                 = "https://console.aws.amazon.com/console/home"
	ConsoleFederationSigninURL = "https://signin.aws.amazon.com/federation"

	ConsoleMinDuration     = 15 * time.Minute
	ConsoleMaxDuration     = 12 * time.Hour
	ConsoleDefaultDuration = 1 * time.Hour
)

type Console struct {
	VaultName string
	Role      string
	Duration  time.Duration
}

type TokenParams struct {
	awsKey   *vaulted.AWSKey
	duration time.Duration
}

func (c *Console) Run(store vaulted.Store) error {
	signinToken, err := c.getSigninToken(store)
	if err != nil {
		return err
	}

	return openConsole(signinToken)
}

func (c *Console) getSigninToken(store vaulted.Store) (string, error) {
	params, err := c.getTokenParams(store)
	if err != nil {
		return "", err
	}

	if c.Role != "" {
		return c.getAssumeRoleToken(store, params)
	} else {
		return c.getFederationToken(params)
	}
}

func (c *Console) getTokenParams(store vaulted.Store) (TokenParams, error) {
	vault, err := c.getVault(store)
	if err != nil {
		return TokenParams{}, err
	}

	awsKey := c.validateAWSKey(vault.AWSKey)

	duration, err := c.chooseDuration(vault.Duration)
	if err != nil {
		return TokenParams{}, err
	}

	params := TokenParams{
		awsKey:   awsKey,
		duration: duration,
	}
	return params, nil
}

func (c *Console) getVault(store vaulted.Store) (*vaulted.Vault, error) {
	vault := &vaulted.Vault{}
	var err error
	if c.VaultName != "" {
		vault, _, err = store.OpenVault(c.VaultName)
		if err != nil {
			return nil, err
		}
	}
	return vault, nil
}

func (c *Console) validateAWSKey(awsKey *vaulted.AWSKey) *vaulted.AWSKey {
	key := &vaulted.AWSKey{}

	if awsKey != nil && awsKey.Valid() {
		key = awsKey
	}

	if c.Role != "" {
		key.Role = c.Role
	}
	return key
}

func (c *Console) chooseDuration(vaultDuration time.Duration) (time.Duration, error) {
	duration := ConsoleDefaultDuration

	if vaultDuration != 0 {
		duration = vaultDuration
	}

	if c.Duration != 0 {
		duration = c.Duration
	}

	return capDuration(duration)
}

func capDuration(duration time.Duration) (time.Duration, error) {
	if duration < ConsoleMinDuration {
		return time.Duration(0), ErrInvalidDuration
	}
	if duration > ConsoleMaxDuration {
		duration = ConsoleMaxDuration
		fmt.Println("Your vault duration is greater than the max console duration.\nCurrent console session duration set to 12 hours.")
	}
	return duration, nil
}

func (c *Console) getAssumeRoleToken(store vaulted.Store, params TokenParams) (string, error) {
	var err error
	awsCreds, err := c.getCredentials(params.awsKey)
	if err != nil {
		return "", err
	}

	if params.awsKey.MFA != "" {
		tokenCode, err := store.Steward().GetMFAToken(c.VaultName)
		if err != nil {
			return "", err
		}
		awsCreds, err = awsCreds.AssumeRoleWithMFA(params.awsKey.MFA, tokenCode, params.awsKey.Role, ConsoleMinDuration)
	} else {
		awsCreds, err = awsCreds.AssumeRole(params.awsKey.Role, ConsoleMinDuration)
	}
	if err != nil {
		return "", err
	}
	return awsCreds.GetSigninToken(&params.duration)
}

func (c *Console) getFederationToken(params TokenParams) (string, error) {
	awsCreds, err := c.getCredentials(params.awsKey)
	if err != nil {
		return "", err
	}

	awsCreds, err = awsCreds.GetFederationToken(params.duration)
	if err != nil {
		return "", err
	}
	return awsCreds.GetSigninToken(nil)
}

func (c *Console) getCredentials(awsKey *vaulted.AWSKey) (*vaulted.AWSCredentials, error) {
	awsCreds, err := awsKey.AWSCredentials.WithLocalDefault()
	if err != nil {
		if err != credentials.ErrNoValidProvidersFoundInChain {
			return nil, err
		} else if err == credentials.ErrNoValidProvidersFoundInChain {
			return nil, ErrNoCredentialsFound
		}
	}

	if awsCreds.ValidSession() {
		return nil, ErrInvalidTemporaryCredentials
	}

	return awsCreds, nil
}

func openConsole(signinToken string) error {
	signinURL, _ := url.Parse(ConsoleFederationSigninURL)
	loginQuery := url.Values{
		"Action":      []string{"login"},
		"SigninToken": []string{signinToken},
		"Destination": []string{ConsoleURL},
	}
	signinURL.RawQuery = loginQuery.Encode()
	err := browser.OpenURL(signinURL.String())
	if err != nil {
		return err
	}
	return nil
}

package vaulted

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/miquella/ask"
)

const (
	DefaultSessionName = "VaultedSession"
)

var STSDurationDefault = time.Hour

var (
	ErrInvalidCommand = errors.New("Invalid command")
)

type Vault struct {
	Duration time.Duration     `json:"duration,omitempty"`
	AWSKey   *AWSKey           `json:"aws_key,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
	SSHKeys  map[string]string `json:"ssh_keys,omitempty"`
}

type AWSCredentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Token  string `json:"token,omitempty"`
}

type AWSKey struct {
	AWSCredentials
	MFA                     string `json:"mfa,omitempty"`
	Role                    string `json:"role,omitempty"`
	ForgoTempCredGeneration bool   `json:"forgoTempCredGeneration"`
}

func (v *Vault) CreateEnvironment(extraVars map[string]string) (*Environment, error) {
	var duration time.Duration
	if v.Duration == 0 {
		duration = STSDurationDefault
	} else {
		duration = v.Duration
	}

	e := &Environment{
		Vars:       make(map[string]string),
		Expiration: time.Now().Add(duration),
	}

	// copy the vault vars to the environment
	for key, value := range v.Vars {
		e.Vars[key] = value
	}
	for key, value := range extraVars {
		e.Vars[key] = value
	}

	// copy the vault ssh keys to the environment
	if len(v.SSHKeys) > 0 {
		e.SSHKeys = make(map[string]string)
		for key, value := range v.SSHKeys {
			e.SSHKeys[key] = value
		}
	}

	// get aws creds
	if v.AWSKey != nil && v.AWSKey.ID != "" && v.AWSKey.Secret != "" {
		if v.AWSKey.ForgoTempCredGeneration {
			e.AWSCreds = &AWSCredentials{
				ID:     v.AWSKey.ID,
				Secret: v.AWSKey.Secret,
			}
		} else {
			var err error
			if v.AWSKey.Role != "" {
				e.AWSCreds, err = v.AWSKey.assumeRole(duration)
			} else {
				e.AWSCreds, err = v.AWSKey.generateSTS(duration)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return e, nil
}

func (k *AWSKey) stsClient() *sts.STS {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			k.ID,
			k.Secret,
			"", // Temporary session token
		),
	})
	return sts.New(sess)
}

func (k *AWSKey) assumeRole(duration time.Duration) (*AWSCredentials, error) {
	// first generate a session token
	creds, err := k.generateSTS(duration)
	if err != nil {
		return nil, err
	}

	// now use the generated session token to assume the role
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			creds.ID,
			creds.Secret,
			creds.Token,
		),
	})

	client := sts.New(sess)
	return k.assumeRoleWithClient(client, duration)
}

func (k *AWSKey) assumeRoleWithClient(client *sts.STS, duration time.Duration) (*AWSCredentials, error) {
	roleSessionName := DefaultSessionName

	callerIdentity, err := client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err == nil {
		parts := strings.SplitN(*callerIdentity.Arn, ":", 6)
		if len(parts) == 6 {
			roleSessionName = fmt.Sprintf("%s@%s", path.Base(parts[5]), parts[4])
		}
	}

	assumeRoleInput := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		RoleArn:         &k.Role,
		RoleSessionName: &roleSessionName,
	}

	assumeRoleOutput, err := client.AssumeRole(assumeRoleInput)
	if err != nil {
		return nil, err
	}

	credentials := &AWSCredentials{
		ID:     *assumeRoleOutput.Credentials.AccessKeyId,
		Secret: *assumeRoleOutput.Credentials.SecretAccessKey,
		Token:  *assumeRoleOutput.Credentials.SessionToken,
	}
	return credentials, nil
}

func (k *AWSKey) buildSessionTokenInput(duration time.Duration) (*sts.GetSessionTokenInput, error) {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	}

	if k.MFA != "" {
		tokenCode, err := getTokenCode()
		if err != nil {
			return nil, err
		}
		input.SerialNumber = &k.MFA
		input.TokenCode = &tokenCode
	}

	return input, nil
}

func (k *AWSKey) generateSTS(duration time.Duration) (*AWSCredentials, error) {
	sessionTokenInput, err := k.buildSessionTokenInput(duration)
	if err != nil {
		return nil, err
	}

	resp, err := k.stsClient().GetSessionToken(sessionTokenInput)
	if err != nil {
		return nil, err
	}

	credentials := &AWSCredentials{
		ID:     *resp.Credentials.AccessKeyId,
		Secret: *resp.Credentials.SecretAccessKey,
		Token:  *resp.Credentials.SessionToken,
	}
	return credentials, nil
}

func getTokenCode() (string, error) {
	tokenCode, err := ask.Ask("Enter your MFA code: ")
	if err != nil {
		return "", err
	}
	tokenCode = strings.TrimSpace(tokenCode)
	return tokenCode, nil
}

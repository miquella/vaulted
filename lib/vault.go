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
		creds, err := v.AWSKey.GetAWSCredentials(duration)
		if err != nil {
			return nil, err
		}

		if v.AWSKey.Role != "" {
			creds, err = creds.AssumeRole(v.AWSKey.Role, duration)
			if err != nil {
				return nil, err
			}
		}

		e.AWSCreds = creds
	}

	return e, nil
}

type AWSKey struct {
	AWSCredentials
	MFA                     string `json:"mfa,omitempty"`
	Role                    string `json:"role,omitempty"`
	ForgoTempCredGeneration bool   `json:"forgoTempCredGeneration"`
}

func (k *AWSKey) GetAWSCredentials(duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	if k.MFA == "" {
		return k.AWSCredentials.GetSessionToken(duration)
	}

	tokenCode, err := getTokenCode()
	if err != nil {
		return nil, err
	}

	return k.AWSCredentials.GetSessionTokenWithMFA(k.MFA, tokenCode, duration)
}

type AWSCredentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Token  string `json:"token,omitempty"`
}

func AWSCredentialsFromSTSCredentials(creds *sts.Credentials) *AWSCredentials {
	return &AWSCredentials{
		ID:     *creds.AccessKeyId,
		Secret: *creds.SecretAccessKey,
		Token:  *creds.SessionToken,
	}
}

func (c *AWSCredentials) GetSessionToken(duration time.Duration) (*AWSCredentials, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}

	getSessionToken, err := client.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(getSessionToken.Credentials), nil
}

func (c *AWSCredentials) GetSessionTokenWithMFA(serialNumber, token string, duration time.Duration) (*AWSCredentials, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}

	getSessionToken, err := client.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(token),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(getSessionToken.Credentials), nil
}

func (c *AWSCredentials) AssumeRole(arn string, duration time.Duration) (*AWSCredentials, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}

	assumeRole, err := client.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(roleSessionName(client)),
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(assumeRole.Credentials), nil
}

func (c *AWSCredentials) client() (*sts.STS, error) {
	// if c is nil, the default credential provider chain is used
	// (yes, I know this seems a little weird)
	config := &aws.Config{}
	if c != nil && c.ID != "" {
		config.Credentials = credentials.NewStaticCredentials(
			c.ID,
			c.Secret,
			c.Token,
		)
	}

	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sts.New(s), nil
}

func getTokenCode() (string, error) {
	tokenCode, err := ask.Ask("Enter your MFA code: ")
	if err != nil {
		return "", err
	}
	tokenCode = strings.TrimSpace(tokenCode)
	return tokenCode, nil
}

func roleSessionName(client *sts.STS) string {
	roleSessionName := DefaultSessionName

	callerIdentity, err := client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err == nil {
		parts := strings.SplitN(*callerIdentity.Arn, ":", 6)
		if len(parts) == 6 {
			roleSessionName = fmt.Sprintf("%s@%s", path.Base(parts[5]), parts[4])
		}
	}

	return roleSessionName
}

package vaulted

import (
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/miquella/ask"
)

var STSDurationDefault = time.Hour

var (
	ErrInvalidCommand = errors.New("Invalid command")
)

type Vault struct {
	Vars     map[string]string `json:"vars"`
	AWSKey   *AWSKey           `json:"aws_key,omitempty"`
	SSHKeys  map[string]string `json:"ssh_keys,omitempty"`
	Duration time.Duration     `json:"duration,omitempty"`
}

type AWSKey struct {
	ID                      string `json:"id"`
	Secret                  string `json:"secret"`
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
		Expiration: time.Now().Add(duration).Unix(),
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
		var err error
		var creds *sts.Credentials
		if v.AWSKey.ForgoTempCredGeneration {
			creds = &sts.Credentials{
				AccessKeyId:     aws.String(v.AWSKey.ID),
				SecretAccessKey: aws.String(v.AWSKey.Secret),
				SessionToken:    aws.String(""), // permanent key, no session token
			}
		} else {
			if v.AWSKey.Role != "" {
				creds, err = v.AWSKey.assumeRole(duration)
			} else {
				creds, err = v.AWSKey.generateSTS(duration)
			}
		}

		if err != nil {
			return nil, err
		}

		e.Vars["AWS_ACCESS_KEY_ID"] = *creds.AccessKeyId
		e.Vars["AWS_SECRET_ACCESS_KEY"] = *creds.SecretAccessKey
		e.Vars["AWS_SESSION_TOKEN"] = *creds.SessionToken
		e.Vars["AWS_SECURITY_TOKEN"] = *creds.SessionToken
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

func (k *AWSKey) getAssumeRoleInput(duration time.Duration) (*sts.AssumeRoleInput, error) {
	roleSessionName := "VaultedSession"
	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		RoleArn:         &k.Role,
		RoleSessionName: &roleSessionName,
	}

	return input, nil
}

func (k *AWSKey) assumeRole(duration time.Duration) (*sts.Credentials, error) {
	// first generate a session token
	creds, err := k.generateSTS(duration)
	if err != nil {
		return nil, err
	}

	// now use the generated session token to assume the role
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			*creds.AccessKeyId,
			*creds.SecretAccessKey,
			*creds.SessionToken,
		),
	})
	stsClient := sts.New(sess)

	assumeRoleInput, err := k.getAssumeRoleInput(duration)
	if err != nil {
		return nil, err
	}

	resp, err := stsClient.AssumeRole(assumeRoleInput)
	if err != nil {
		return nil, err
	}

	return resp.Credentials, nil
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

func (k *AWSKey) generateSTS(duration time.Duration) (*sts.Credentials, error) {
	sessionTokenInput, err := k.buildSessionTokenInput(duration)
	if err != nil {
		return nil, err
	}

	resp, err := k.stsClient().GetSessionToken(sessionTokenInput)
	if err != nil {
		return nil, err
	}

	return resp.Credentials, nil
}

func getTokenCode() (string, error) {
	tokenCode, err := ask.Ask("Enter your MFA code: ")
	if err != nil {
		return "", err
	}
	tokenCode = strings.TrimSpace(tokenCode)
	return tokenCode, nil
}

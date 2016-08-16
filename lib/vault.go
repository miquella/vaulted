package vaulted

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
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
	ID     string `json:"id"`
	Secret string `json:"secret"`
	MFA    string `json:"mfa,omitempty"`
	Role   string `json:"role,omitempty"`
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
		var stsCreds map[string]string
		if v.AWSKey.Role != "" {
			stsCreds, err = v.AWSKey.assumeRole(duration)
		} else {
			stsCreds, err = v.AWSKey.generateSTS(duration)
		}
		if err != nil {
			return nil, err
		}

		for key, value := range stsCreds {
			e.Vars[key] = value
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

func (k *AWSKey) getAssumeRoleInput(duration time.Duration) *sts.AssumeRoleInput {
	roleSessionName := "VaultedSession"
	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		RoleArn:         &k.Role,
		RoleSessionName: &roleSessionName,
	}

	if k.MFA != "" {
		tokenCode := getTokenCode()
		input.SerialNumber = &k.MFA
		input.TokenCode = &tokenCode
	}

	return input
}

func (k *AWSKey) assumeRole(duration time.Duration) (map[string]string, error) {
	resp, err := k.stsClient().AssumeRole(k.getAssumeRoleInput(duration))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"AWS_ACCESS_KEY_ID":     *resp.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY": *resp.Credentials.SecretAccessKey,
		"AWS_SESSION_TOKEN":     *resp.Credentials.SessionToken,
	}, nil
}

func (k *AWSKey) buildSessionTokenInput(duration time.Duration) *sts.GetSessionTokenInput {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	}

	if k.MFA != "" {
		tokenCode := getTokenCode()
		input.SerialNumber = &k.MFA
		input.TokenCode = &tokenCode
	}

	return input
}

func (k *AWSKey) generateSTS(duration time.Duration) (map[string]string, error) {
	resp, err := k.stsClient().GetSessionToken(k.buildSessionTokenInput(duration))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"AWS_ACCESS_KEY_ID":     *resp.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY": *resp.Credentials.SecretAccessKey,
		"AWS_SESSION_TOKEN":     *resp.Credentials.SessionToken,
	}, nil
}

func getTokenCode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your MFA code: ")
	tokenCode, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	tokenCode = strings.TrimSpace(tokenCode)
	return tokenCode
}

package vaulted

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	STS_DURATION = 3600
)

var (
	ErrInvalidCommand = errors.New("Invalid command")
)

type Vault struct {
	Vars    map[string]string `json:"vars"`
	AWSKey  *AWSKey           `json:"aws_key,omitempty"`
	SSHKeys map[string]string `json:"ssh_keys,omitempty"`
}

type AWSKey struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	MFA    string `json:"mfa,omitempty"`
	Role   string `json:"role,omitempty"`
}

func (v *Vault) CreateEnvironment(extraVars map[string]string) (*Environment, error) {
	e := &Environment{
		Vars:       make(map[string]string),
		Expiration: time.Now().Add(STS_DURATION * time.Second).Unix(),
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
			stsCreds, err = v.AWSKey.assumeRole()
		} else {
			stsCreds, err = v.AWSKey.generateSTS()
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

func (v *Vault) GetEnvironment(name, password string) (*Environment, error) {
	env, err := v.openEnvironment(name, password)
	if err == nil {
		expired := time.Now().Add(5 * time.Minute).After(time.Unix(env.Expiration, 0))
		if !expired {
			return env, nil
		}
	}

	// the environment isn't valid (possibly expired), so remove it
	removeEnvironment(name)

	env, err = v.CreateEnvironment(map[string]string{"VAULTED_ENV": name})
	if err != nil {
		return nil, err
	}

	// we have a valid environment, so if saving fails, ignore the failure
	v.sealEnvironment(name, password, env)
	return env, nil
}

func (v *Vault) sealEnvironment(name, password string, env *Environment) error {
	// read the vault file (to get key details)
	vf, err := readVaultFile(name)
	if err != nil {
		return err
	}

	// marshal the environment content
	content, err := json.Marshal(env)
	if err != nil {
		return err
	}

	// encrypt the environment
	ef := &EnvironmentFile{
		Method:  "secretbox",
		Details: make(Details),
	}

	switch ef.Method {
	case "secretbox":
		nonce := [24]byte{}
		_, err = rand.Read(nonce[:])
		if err != nil {
			return err
		}
		ef.Details.SetBytes("nonce", nonce[:])

		key := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(key))
		if err != nil {
			return err
		}
		copy(key[:], derivedKey[:])

		ef.Ciphertext = secretbox.Seal(nil, content, &nonce, &key)

	default:
		return err
	}

	return writeEnvironmentFile(name, ef)
}

func (v *Vault) openEnvironment(name, password string) (*Environment, error) {
	vf, err := readVaultFile(name)
	if err != nil {
		return nil, err
	}

	ef, err := readEnvironmentFile(name)
	if err != nil {
		return nil, err
	}

	e := Environment{}

	switch ef.Method {
	case "secretbox":
		if vf.Key == nil {
			return nil, ErrInvalidKeyConfig
		}

		nonce := ef.Details.Bytes("nonce")
		if len(nonce) == 0 {
			return nil, ErrInvalidEncryptionConfig
		}
		boxNonce := [24]byte{}
		copy(boxNonce[:], nonce)

		boxKey := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(boxKey))
		if err != nil {
			return nil, err
		}
		copy(boxKey[:], derivedKey[:])

		plaintext, ok := secretbox.Open(nil, ef.Ciphertext, &boxNonce, &boxKey)
		if !ok {
			return nil, ErrInvalidPassword
		}

		err = json.Unmarshal(plaintext, &e)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid encryption method: %s", ef.Method)
	}

	return &e, nil
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

func (k *AWSKey) getAssumeRoleInput() *sts.AssumeRoleInput {
	roleSessionName := "VaultedSession"
	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(STS_DURATION),
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

func (k *AWSKey) assumeRole() (map[string]string, error) {
	resp, err := k.stsClient().AssumeRole(k.getAssumeRoleInput())
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"AWS_ACCESS_KEY_ID":     *resp.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY": *resp.Credentials.SecretAccessKey,
		"AWS_SESSION_TOKEN":     *resp.Credentials.SessionToken,
	}, nil
}

func (k *AWSKey) buildSessionTokenInput() *sts.GetSessionTokenInput {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(STS_DURATION),
	}

	if k.MFA != "" {
		tokenCode := getTokenCode()
		input.SerialNumber = &k.MFA
		input.TokenCode = &tokenCode
	}

	return input
}

func (k *AWSKey) generateSTS() (map[string]string, error) {
	resp, err := k.stsClient().GetSessionToken(k.buildSessionTokenInput())
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

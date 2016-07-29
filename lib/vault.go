package vaulted

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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

func (v *Vault) CreateEnvironment(staticEnvironment bool, extraVars map[string]string) (map[string]string, error) {
	vars := make(map[string]string)
	for key, value := range v.Vars {
		vars[key] = value
	}
	for key, value := range extraVars {
		vars[key] = value
	}

	// start ssh agent in dynamic environments
	if !staticEnvironment {
		sock, err := v.startProxyKeyring()
		if err != nil {
			return nil, err
		}

		vars["SSH_AUTH_SOCK"] = sock
	}

	// get aws creds (use sts in dynamic environments)
	if v.AWSKey != nil && v.AWSKey.ID != "" && v.AWSKey.Secret != "" {
		var err error
		var stsCreds map[string]string
		if staticEnvironment {
			stsCreds["AWS_ACCESS_KEY_ID"] = v.AWSKey.ID
			stsCreds["AWS_SECRET_ACCESS_KEY"] = v.AWSKey.Secret
		} else {
			if v.AWSKey.Role != "" {
				stsCreds, err = v.AWSKey.assumeRole()
			} else {
				stsCreds, err = v.AWSKey.generateSTS()
			}
		}
		if err != nil {
			return nil, err
		}

		for key, value := range stsCreds {
			vars[key] = value
		}
	}

	return vars, nil
}

func (v *Vault) Spawn(cmd []string, extraVars map[string]string) (*int, error) {
	if len(cmd) == 0 {
		return nil, ErrInvalidCommand
	}

	// lookup the path of the executable
	cmdpath, err := exec.LookPath(cmd[0])
	if err != nil {
		return nil, fmt.Errorf("Cannot find executable %s: %v", cmd[0], err)
	}

	// build the environ
	vars, err := v.CreateEnvironment(false, extraVars)
	if err != nil {
		return nil, err
	}

	// start the process
	var attr os.ProcAttr
	attr.Env = v.getEnviron(vars)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	proc, err := os.StartProcess(cmdpath, cmd, &attr)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute command: %v", err)
	}

	// wait for the process to exit
	state, _ := proc.Wait()

	var exitStatus int
	if !state.Success() {
		if status, ok := state.Sys().(syscall.WaitStatus); ok {
			exitStatus = status.ExitStatus()
		} else {
			exitStatus = 255
		}
	}

	// we only return an error if spawning the process failed, not if
	// the spawned command returned a failure status code.
	return &exitStatus, nil
}

func (v *Vault) getEnviron(vars map[string]string) []string {
	// load the current environ
	env := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		env[parts[0]] = parts[1]
	}

	// merge the provided vars
	for key, value := range vars {
		env[key] = value
	}

	// recombine into environ
	environ := make([]string, 0, len(env))
	for key, value := range env {
		environ = append(environ, fmt.Sprintf("%s=%s", key, value))
	}
	return environ
}

func (v *Vault) startProxyKeyring() (string, error) {
	keyring, err := NewProxyKeyring(os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return "", err
	}

	// load ssh keys
	for comment, key := range v.SSHKeys {
		addedKey := agent.AddedKey{
			Comment: comment,
		}

		addedKey.PrivateKey, err = ssh.ParseRawPrivateKey([]byte(key))
		if err != nil {
			return "", err
		}

		err := keyring.Add(addedKey)
		if err != nil {
			return "", err
		}
	}

	sock, err := keyring.Listen()
	if err != nil {
		return "", err
	}

	go keyring.Serve()

	return sock, err
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

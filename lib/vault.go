package vaulted

import (
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
)

const (
	STS_DURATION = 3600
)

var (
	ErrInvalidCommand = errors.New("Invalid command")
)

type Vault struct {
	Vars   map[string]string `json:"vars"`
	AWSKey *AWSKey           `json:"aws_key,omitempty"`
}

type AWSKey struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func (v *Vault) Spawn(cmd []string, env map[string]string) (*int, error) {
	if len(cmd) == 0 {
		return nil, ErrInvalidCommand
	}

	// lookup the path of the executable
	cmdpath, err := exec.LookPath(cmd[0])
	if err != nil {
		return nil, fmt.Errorf("Cannot find executable %s: %v", cmd[0], err)
	}

	// build the environ
	vars := make(map[string]string)
	for key, value := range v.Vars {
		vars[key] = value
	}
	for key, value := range env {
		vars[key] = value
	}

	if v.AWSKey != nil && v.AWSKey.ID != "" && v.AWSKey.Secret != "" {
		stsCreds, err := v.AWSKey.generateSTS()
		if err != nil {
			return nil, err
		}
		for k, v := range stsCreds {
			vars[k] = v
		}
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

func (k *AWSKey) generateSTS() (map[string]string, error) {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			k.ID,
			k.Secret,
			"", // Temporary session token
		),
	})

	params := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(STS_DURATION),
	}

	stsClient := sts.New(sess)
	resp, err := stsClient.GetSessionToken(params)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"AWS_ACCESS_KEY_ID":     *resp.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY": *resp.Credentials.SecretAccessKey,
		"AWS_SESSION_TOKEN":     *resp.Credentials.SessionToken,
	}, nil
}

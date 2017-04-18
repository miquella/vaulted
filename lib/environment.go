package vaulted

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Environment struct {
	Expiration time.Time         `json:"expiration"`
	AWSCreds   *AWSCredentials   `json:"aws_creds,omitempty"`
	Vars       map[string]string `json:"vars,omitempty"`
	SSHKeys    map[string]string `json:"ssh_keys,omitempty"`
}

func (e *Environment) Spawn(cmd []string, extraVars map[string]string) (*int, error) {
	if len(cmd) == 0 {
		return nil, ErrInvalidCommand
	}

	// lookup the path of the executable
	cmdpath, err := exec.LookPath(cmd[0])
	if err != nil {
		return nil, fmt.Errorf("Cannot find executable %s: %v", cmd[0], err)
	}

	// copy the extra vars so we can mutate it
	vars := make(map[string]string)
	for key, value := range extraVars {
		vars[key] = value
	}

	// start the agent
	sock, err := e.startProxyKeyring()
	if err != nil {
		return nil, err
	}

	vars["SSH_AUTH_SOCK"] = sock

	// trap signals
	sigs := make(chan os.Signal)
	signal.Notify(
		sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGWINCH,
	)

	// start the process
	var attr os.ProcAttr
	attr.Env = e.buildEnviron(vars)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	proc, err := os.StartProcess(cmdpath, cmd, &attr)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute command: %v", err)
	}

	// relay trapped signals to the spawned process
	go func() {
		for s := range sigs {
			proc.Signal(s)
		}
	}()

	defer func() {
		signal.Stop(sigs)
		close(sigs)
	}()

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

func (e *Environment) startProxyKeyring() (string, error) {
	keyring, err := NewProxyKeyring(os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return "", err
	}

	// load ssh keys
	for comment, key := range e.SSHKeys {
		timeRemaining := e.Expiration.Sub(time.Now())
		addedKey := agent.AddedKey{
			Comment:      comment,
			LifetimeSecs: uint32(timeRemaining.Seconds()),
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

func (e *Environment) Variables() *Variables {
	vars := Variables{
		Set: make(map[string]string),
	}

	for key, value := range e.Vars {
		vars.Set[key] = value
	}

	if e.AWSCreds != nil {
		vars.Set["AWS_ACCESS_KEY_ID"] = e.AWSCreds.ID
		vars.Set["AWS_SECRET_ACCESS_KEY"] = e.AWSCreds.Secret

		if e.AWSCreds.Token != "" {
			vars.Set["AWS_SESSION_TOKEN"] = e.AWSCreds.Token
			vars.Set["AWS_SECURITY_TOKEN"] = e.AWSCreds.Token
		} else {
			vars.Unset = append(
				vars.Unset,
				"AWS_SESSION_TOKEN",
				"AWS_SECURITY_TOKEN",
			)
		}
	}

	return &vars
}

func (e *Environment) buildEnviron(extraVars map[string]string) []string {
	vars := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		vars[parts[0]] = parts[1]
	}

	v := e.Variables()
	for _, key := range v.Unset {
		delete(vars, key)
	}
	for key, value := range v.Set {
		vars[key] = value
	}

	for key, value := range extraVars {
		vars[key] = value
	}

	environ := make([]string, 0, len(vars))
	for key, value := range vars {
		environ = append(environ, fmt.Sprintf("%s=%s", key, value))
	}
	return environ
}

type Variables struct {
	Set   map[string]string
	Unset []string
}

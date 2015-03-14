package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/miquella/vaulted/vault"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	filename string
	password string = "password"
	envs     Envs   = Envs{}

	accountJson  []byte
	accountVault *vault.AccountVault = &vault.AccountVault{}
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&filename, "f", filepath.Join(u.HomeDir, ".shell-vault"), "vault filename")
	flag.Var(envs, "e", "env vars")
	flag.Parse()

	if flag.NArg() < 1 {
		println("invalid syntax")
		os.Exit(1)
	}
}

func main() {
	var err error
	if _, err = os.Stat(filename); !os.IsNotExist(err) {
		accountVault, err = vault.LoadAccountVault(filename, password)
		if err != nil {
			panic(err)
		}

		accountJson, err = json.Marshal(accountVault)
		if err != nil {
			panic(err)
		}
	}

	switch flag.Arg(0) {
	case "add":
		if flag.NArg() != 2 {
			println("invalid syntax")
			os.Exit(1)
		}
		addAccount(flag.Arg(1))
	case "cat":
		if flag.NArg() != 2 {
			println("invalid syntax")
			os.Exit(1)
		}
		showAccount(flag.Arg(1))
	case "list":
		if flag.NArg() != 1 {
			println("invalid syntax")
			os.Exit(1)
		}
		listAccounts()
	case "shell":
		if flag.NArg() != 2 {
			println("invalid syntax")
			os.Exit(1)
		}
		shellAccount(flag.Arg(1))
	default:
		println("unknown subcommand")
		os.Exit(1)
	}

	// save vault if it has changed
	newAccountJson, err := json.Marshal(accountVault)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(accountJson, newAccountJson) {
		err = accountVault.SaveAccountVault(filename, password)
		if err != nil {
			println("error updating vault")
			println(err)
			os.Exit(2)
		}
		println("vault updated")
	}
}

func addAccount(name string) {
	if _, exists := accountVault.Accounts[name]; exists {
		println("ERROR: account already exists")
		os.Exit(2)
	}

	if accountVault.Accounts == nil {
		accountVault.Accounts = make(map[string]vault.Account)
	}

	account := vault.Account{
		Name: name,
		Env:  map[string]string{},
	}
	for v, val := range envs {
		account.Env[v] = val
	}
	accountVault.Accounts[name] = account
}

func showAccount(name string) {
	account, exists := accountVault.Accounts[name]
	if !exists {
		println("unknown account name")
		os.Exit(2)
	}

	for k, val := range account.Env {
		fmt.Printf("%s=%s\n", k, val)
	}
}

func listAccounts() {
	for name := range accountVault.Accounts {
		fmt.Println(name)
	}
}

func shellAccount(name string) {
	account, exists := accountVault.Accounts[name]
	if !exists {
		println("unknown account name")
		os.Exit(2)
	}

	envs := ParseEnviron(os.Environ())
	for key, val := range account.Env {
		envs[key] = val
	}

	cmd := exec.Command(envs["SHELL"], "--login")
	cmd.Env = CreateEnviron(envs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	if cmd.ProcessState.Success() {
		os.Exit(0)
	} else {
		os.Exit(5)
	}
}

type Envs map[string]string

func (e Envs) String() string {
	return "[]"
}

func (e Envs) Set(v string) error {
	parts := strings.SplitN(v, "=", 2)
	if len(parts) < 2 {
		return errors.New("An env var must be in the form VAR=VALUE")
	}

	if _, exists := e[parts[0]]; exists {
		return errors.New("An env var cannot be set more than once")
	}

	e[parts[0]] = parts[1]
	return nil
}

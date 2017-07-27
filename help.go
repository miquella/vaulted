package main

//go:generate sh -c "for md in doc/vaulted*.md; do md2man-roff ${DOLLAR}md > doc/man/${DOLLAR}(basename ${DOLLAR}md .md); done"
//go:generate go-bindata -nometadata -prefix doc/man/ -o man.go doc/man/
//go:generate gofmt -w man.go

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/miquella/vaulted/lib"
)

var (
	ErrHelp = errors.New("help requested")

	HelpAliases = map[string]string{
		"add":     "add",
		"create":  "add",
		"new":     "add",
		"cp":      "cp",
		"copy":    "cp",
		"dump":    "dump",
		"edit":    "edit",
		"env":     "env",
		"ls":      "ls",
		"list":    "ls",
		"load":    "load",
		"rm":      "rm",
		"delete":  "rm",
		"remove":  "rm",
		"shell":   "shell",
		"upgrade": "upgrade",
	}
)

type Help struct {
	Subcommand string
}

func (h *Help) Run(store vaulted.Store) error {
	return displayHelp(h.Subcommand)
}

func displayHelp(subcommand string) error {
	if subcommand != "" {
		if HelpAliases[subcommand] == "" {
			return fmt.Errorf("Unknown subcommand: %s", subcommand)
		}
		subcommand = HelpAliases[subcommand]
	}

	manpage := "vaulted.1"
	if subcommand != "" {
		manpage = fmt.Sprintf("vaulted-%s.1", subcommand)
	}

	content, err := Asset(manpage)
	if err != nil {
		return err
	}

	manpath, err := exec.LookPath("man")
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "vaulted")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	man := path.Join(dir, manpage)
	err = ioutil.WriteFile(man, content, 0644)
	if err != nil {
		return err
	}

	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	proc, err := os.StartProcess(manpath, []string{"man", man}, &attr)
	if err != nil {
		return err
	}

	state, _ := proc.Wait()
	if !state.Success() {
		os.Exit(255)
	}
	os.Exit(1)
	return nil
}

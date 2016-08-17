package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Env struct {
	VaultName string
	Shell     string
	UsageHint bool
}

func (e *Env) Run(steward Steward) error {
	_, env, err := steward.GetEnvironment(e.VaultName, nil)
	if err != nil {
		return err
	}

	usageHint := ""
	setVar := ""
	quoteReplacement := "\""
	switch e.Shell {
	case "fish":
		usageHint = "# To load these variables into your shell, execute:\n#   eval (%s)\n"
		setVar = "set -x %s \"%s\";\n"
		quoteReplacement = "\\\""
	default:
		usageHint = "# To load these variables into your shell, execute:\n#   eval $(%s)\n"
		setVar = "export %s=\"%s\"\n"
		quoteReplacement = "\\\""
	}

	// sort the vars
	var keys []string
	for key, _ := range env.Vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// display the vars using the format string for the shell
	if e.UsageHint {
		fmt.Printf(usageHint, strings.Join(os.Args, " "))
	}

	for _, key := range keys {
		fmt.Printf(setVar, key, strings.Replace(env.Vars[key], "\"", quoteReplacement, -1))
	}

	return nil
}

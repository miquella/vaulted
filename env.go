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
	unsetVar := ""
	setVar := ""
	quoteReplacement := "\""
	switch e.Shell {
	case "fish":
		usageHint = "# To load these variables into your shell, execute:\n#   eval (%s)\n"
		unsetVar = "set -e %s;\n"
		setVar = "set -x %s \"%s\";\n"
		quoteReplacement = "\\\""
	default:
		usageHint = "# To load these variables into your shell, execute:\n#   eval $(%s)\n"
		unsetVar = "unset %s\n"
		setVar = "export %s=\"%s\"\n"
		quoteReplacement = "\\\""
	}

	vars := env.Variables()

	// sort the vars to unset
	sort.Strings(vars.Unset)

	// sort the vars to set
	var keys []string
	for key := range vars.Set {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// display the vars using the format string for the shell
	if e.UsageHint {
		fmt.Printf(usageHint, strings.Join(os.Args, " "))
	}

	for _, key := range vars.Unset {
		fmt.Printf(unsetVar, key)
	}

	for _, key := range keys {
		fmt.Printf(setVar, key, strings.Replace(vars.Set[key], "\"", quoteReplacement, -1))
	}

	return nil
}

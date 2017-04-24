package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/miquella/vaulted/lib"
)

type Env struct {
	VaultName string
	Role      string

	DetectedShell string
	Format        string
	Command       string
}

type templateVals struct {
	vaulted.Variables
	Command string
}

var (
	envFormatters = map[string]string{
		"fish": `# To load these variables into your shell, execute:
#   eval ({{ .Command }})
{{ range $var := .Unset}}set -e {{ $var }};
{{ end -}}
{{ range $var, $value := .Set }}set -x {{ $var }} "{{ replace $value "\"" "\\\"" }}";
{{ end }}`,
		"sh": `# To load these variables into your shell, execute:
#   eval $({{ .Command }})
{{ range $var := .Unset}}unset {{ $var }}
{{ end -}}
{{ range $var, $value := .Set }}export {{ $var }}="{{ replace $value "\"" "\\\"" }}"
{{ end }}`,
		"json": "{{ json .Set }}\n",
	}
)

var templateFuncMap = template.FuncMap{
	"replace": func(val string, toReplace string, replacement string) string {
		return strings.Replace(val, toReplace, replacement, -1)
	},
	"json": func(val interface{}) (string, error) {
		json, err := json.MarshalIndent(val, "", "  ")
		if err != nil {
			return "", err
		} else {
			return string(json), nil
		}
	},
}

func (e *Env) Run(steward Steward) error {
	env, err := e.getEnvironment(steward)
	if err != nil {
		return err
	}

	var templateStr string
	format := e.Format

	if format == "shell" {
		format = e.DetectedShell
	}

	if foundTemplate, ok := envFormatters[format]; ok {
		templateStr = foundTemplate
	} else {
		templateStr = format
	}
	tmpl, err := template.New("envTmpl").Funcs(templateFuncMap).Parse(templateStr)

	vals := templateVals{}
	variables := env.Variables()

	vals.Set = variables.Set

	sort.Strings(variables.Unset)
	vals.Unset = variables.Unset

	vals.Command = e.Command

	if err != nil {
		return ErrorWithExitCode{err, 64}
	}
	return tmpl.Execute(os.Stdout, vals)
}

func (e *Env) getEnvironment(steward Steward) (*vaulted.Environment, error) {
	var err error

	// default environment
	env := DefaultEnvironment()

	if e.VaultName != "" {
		// get specific environment
		_, env, err = steward.GetEnvironment(e.VaultName, nil)
		if err != nil {
			return nil, err
		}
	}

	if e.Role != "" {
		return env.Assume(e.Role)
	}

	return env, nil
}

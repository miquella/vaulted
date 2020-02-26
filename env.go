package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/miquella/vaulted/v3/lib"
)

type Env struct {
	SessionOptions

	DetectedShell string
	Format        string
	Command       string
	Interactive   bool
}

type templateVals struct {
	Command     string
	Interactive bool

	AWSCreds struct {
		ID     string
		Secret string
		Token  string `json:",omitempty"`
	}

	Set   map[string]string
	Unset []string
}

var (
	sessionFormatters = map[string]string{
		"fish": `
{{- if .Interactive -}}
# To load these variables into your shell, execute:
#   {{ .Command }} | source
{{ end -}}
{{ range $var := .Unset}}set -e {{ $var }};
{{ end -}}
{{ range $var, $value := .Set }}set -gx {{ $var }} "{{ replace $value "\"" "\\\"" }}";
{{ end -}}
`,
		"sh": `
{{- if .Interactive -}}
# To load these variables into your shell, execute:
#   eval "$({{ .Command }})"
{{ end -}}
{{ range $var := .Unset}}unset {{ $var }}
{{ end -}}
{{ range $var, $value := .Set }}export {{ $var }}="{{ replace $value "\"" "\\\"" }}"
{{ end -}}
`,
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

func (e *Env) Run(store vaulted.Store) error {
	session, err := GetSessionWithOptions(store, &e.SessionOptions)
	if err != nil {
		return err
	}

	var templateStr string
	format := e.Format

	if format == "shell" {
		format = e.DetectedShell
		if _, ok := sessionFormatters[format]; !ok {
			format = "sh"
		}
	}

	if foundTemplate, ok := sessionFormatters[format]; ok {
		templateStr = foundTemplate
	} else {
		templateStr = format
	}
	tmpl, err := template.New("sessionTmpl").Funcs(templateFuncMap).Parse(templateStr)

	variables := session.Variables()
	sort.Strings(variables.Unset)

	vals := templateVals{
		Command:     e.Command,
		Interactive: e.Interactive,
		Set:         variables.Set,
		Unset:       variables.Unset,
	}

	if session.AWSCreds != nil {
		vals.AWSCreds.ID = session.AWSCreds.ID
		vals.AWSCreds.Secret = session.AWSCreds.Secret
		vals.AWSCreds.Token = session.AWSCreds.Token
	}

	if err != nil {
		return ErrorWithExitCode{err, 64}
	}
	return tmpl.Execute(os.Stdout, vals)
}

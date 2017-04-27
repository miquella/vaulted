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
	Command string

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
		"fish": `# To load these variables into your shell, execute:
#   eval ({{ .Command }})
{{ range $var := .Unset}}set -e {{ $var }};
{{ end -}}
{{ range $var, $value := .Set }}set -gx {{ $var }} "{{ replace $value "\"" "\\\"" }}";
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
	session, err := e.getSession(steward)
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
		Command: e.Command,
		Set:     variables.Set,
		Unset:   variables.Unset,
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

func (e *Env) getSession(steward Steward) (*vaulted.Session, error) {
	var err error

	// default session
	session := DefaultSession()

	if e.VaultName != "" {
		// get specific session
		_, session, err = steward.GetSession(e.VaultName, nil)
		if err != nil {
			return nil, err
		}
	}

	if e.Role != "" {
		return session.Assume(e.Role)
	}

	return session, nil
}

package clic

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TmplData is the structure used for usage output templating. Custom template
// string values should be based on this type.
type TmplData struct {
	ResolvedCmd    *Clic
	ResolvedCmdSet []*Clic
}

// TmplConfig tracks the template string and function map used for usage output
// templating.
type TmplConfig struct {
	Text string
	FMap template.FuncMap
}

// NewDefaultTmplConfig returns the default TmplConfig value. This can be used
// as an example of how to setup custom usage output templating.
func NewDefaultTmplConfig() *TmplConfig {
	resolvedCmdSetHintFn := func(cmds []*Clic) string {
		var out, sep string
		for _, cmd := range cmds {
			out += sep + cmd.FlagSet.Name()
			sep = " "
			if len(cmd.FlagSet.Flags()) > 0 {
				out += sep + "[FLAGS]"
			}
		}
		return out
	}

	subsAndOperandsHintFn := func(cmd *Clic) string {
		var out, sep string
		var anySubShowing bool

		for _, sub := range cmd.subs {
			if sub.HideUsage {
				continue
			}
			anySubShowing = true

			out += sep + sub.FlagSet.Name()
			sep = "|"
		}

		if anySubShowing {
			pre, suf := "[", "]"
			if cmd.SubRequired {
				pre, suf = "{", "}"
			}
			out = pre + out + suf

			if len(cmd.OperandSet.Operands()) == 0 {
				return " " + out
			}

			out += " | "
			sep = ""
		}

		for _, op := range cmd.OperandSet.Operands() {
			pre, suf := "[", "]"
			if op.IsRequired() {
				pre, suf = "<", ">"
			}
			out += sep + pre + op.Name() + suf
			sep = " "
		}

		pre, suf := "{", "}"
		if !anySubShowing {
			pre, suf = "", ""
		}
		return " " + pre + out + suf
	}

	tmplFMap := template.FuncMap{
		"ResolvedCmdSetHint":  resolvedCmdSetHintFn,
		"SubsAndOperandsHint": subsAndOperandsHintFn,
		"stringsJoin":         strings.Join,
	}

	tmplText := strings.TrimSpace(`
{{- $cmd := .ResolvedCmd -}}
Usage:

{{if .}}  {{end}}{{ResolvedCmdSetHint .ResolvedCmdSet}}
{{- SubsAndOperandsHint $cmd}}
    {{- if $cmd.Description}}

      {{$cmd.Description}}
    {{- end}}
{{if $cmd.FlagSet.Flags}}
{{$cmd.FlagSet.Usage}}{{- end}}
{{- if $cmd.Aliases}}
Aliases for {{$cmd.FlagSet.Name}}:

      {{stringsJoin $cmd.Aliases ", "}}
{{- end}}
`)

	return &TmplConfig{
		Text: tmplText,
		FMap: tmplFMap,
	}
}

func executeTmpl(tc *TmplConfig, data any) string {
	tmpl := template.New("clic").Funcs(tc.FMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(tc.Text)
	if err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, data); err != nil {
		fmt.Fprintf(buf, "%v\n", err)
		return buf.String()
	}

	return buf.String()
}

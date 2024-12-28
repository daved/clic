package operandset

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type TmplData struct {
	OperandSet *OperandSet
}

type TmplConfig struct {
	Text string
	FMap template.FuncMap
}

func NewDefaultTmplConfig() *TmplConfig {
	tagHintFn := func(o *Operand) string {
		if o.Tag == "" {
			return ""
		}

		pre, post := "[", "]"
		if o.req {
			pre, post = "<", ">"
		}

		return pre + o.Tag + post
	}
	nameHintFn := func(o *Operand) string {
		if o.name == "" {
			return ""
		}

		var post string
		if o.req {
			post = "  (required)"
		}

		return o.name + post
	}

	tmplFMap := template.FuncMap{
		"Join":     strings.Join,
		"TagHint":  tagHintFn,
		"NameHint": nameHintFn,
	}

	tmplText := strings.TrimSpace(`
{{- if .OperandSet.Operands -}}
Operands for {{.OperandSet.Name}}:
{{range $i, $op := .OperandSet.Operands}}
  {{if .}}  {{end}}{{if $op.Name}}{{NameHint $op}}{{end}}
        {{$op.Description}}
{{end}}{{else}}{{- end}}
`)

	return &TmplConfig{
		Text: tmplText,
		FMap: tmplFMap,
	}
}

func executeTmpl(tc *TmplConfig, data any) string {
	tmpl := template.New("operandset").Funcs(tc.FMap)

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

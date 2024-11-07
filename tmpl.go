package clic

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

type tmplData struct {
	Root    *Clic
	Current *Clic
	Called  []*Clic
}

var tmplText = strings.TrimSpace(`
{{- $cur := .Current -}}{{- $subsStarted := false -}}
{{- $leftBrack := "[" -}}{{- $rightBrack := "]" -}}
Usage:

{{if .}}  {{end}}{{range $clic := .Called}}
  {{- $clic.Handler.FlagSet.Name}} {{if $clic.Handler.FlagSet.Opts}}[FLAGS] {{end -}}
  {{- if eq $cur.Handler.FlagSet.Name $clic.Handler.FlagSet.Name }}
    {{- if $cur.Meta.SubRequired}}{{$leftBrack = "{"}}{{$rightBrack = "}"}}{{end -}}
    {{- if $cur.Subs}}{{$leftBrack}}{{end}}{{range $i, $sub := $cur.Subs}}
      {{- if $sub.Meta.SkipUsage}}{{continue}}{{end}}
      {{- if and $i $subsStarted}}|{{end}}{{$sub.Handler.FlagSet.Name}}{{$subsStarted = true}}
    {{- end}}{{if $cur.Subs}}{{$rightBrack}}{{end}}
    {{- if $clic.Meta.ArgsHint}}{{$clic.Meta.ArgsHint}}{{end}}
    {{- if $clic.Meta.CmdDesc}}

      {{$clic.Meta.CmdDesc}}
    {{- end}}
  {{- end}}
{{- end}}

{{.Current.Handler.FlagSet.Usage}}
`)

func (c *Clic) Usage() string {
	data := &tmplData{
		Root:    root(c),
		Current: c,
		Called:  allCalled(c),
	}

	tmpl := template.New("clic")

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(c.tmplTxt)
	if err != nil {
		fmt.Fprintf(buf, "cli command: template error: %v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, data); err != nil {
		fmt.Fprintf(buf, "cli command: template error: %v\n", err)
		return buf.String()
	}

	return buf.String()
}

func root(c *Clic) *Clic {
	root := c

	for root.Parent != nil {
		root = root.Parent
	}

	return root
}

func allCalled(c *Clic) []*Clic {
	var all []*Clic
	cur := c

	for cur.Parent != nil {
		all = append(all, cur)
		cur = cur.Parent
	}
	all = append(all, cur)

	slices.Reverse(all)
	return all
}

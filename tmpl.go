package clic

import (
	"strings"
)

type tmplData struct {
	Current *Clic
	Called  []*Clic
}

var tmplText = strings.TrimSpace(`
{{- $cur := .Current -}}
{{- $leftBrack := "[" -}}{{- $rightBrack := "]" -}}
Usage:

{{if .}}  {{end}}{{range $clic := .Called}}
  {{- $clic.FlagSet.Name}} {{if $clic.FlagSet.Flags}}[FLAGS] {{end -}}
  {{- if eq $cur.FlagSet.Name $clic.FlagSet.Name }}
    {{- if $cur.SubRequired}}{{$leftBrack = "{"}}{{$rightBrack = "}"}}{{end -}}
    {{- if $cur.Subs}}{{$leftBrack}}{{end}}{{range $i, $sub := $cur.Subs}}
      {{- if $sub.UsageConfig.Skip}}{{continue}}{{end}}
      {{- if $i}}|{{end}}{{$sub.FlagSet.Name}}
    {{- end}}{{/* range sub */}}
    {{- if and $cur.Subs $cur.OperandSet.Operands}}|{{end}}
    {{- range $i, $op := $cur.OperandSet.Operands}}{{if $i}} {{end}}{{$op.Tag}}{{end}}
    {{- if $cur.Subs}}{{$rightBrack}}{{end}}
    {{- if $clic.UsageConfig.CmdDesc}}

      {{$clic.UsageConfig.CmdDesc}}
    {{- end}}{{/* CmdDesc */}}
  {{- end}}{{/* eq Name Name */}}
{{- end}}{{/* range clic */}}

{{.Current.FlagSet.Usage}}
`)

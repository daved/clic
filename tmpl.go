package clic

import (
	"strings"
)

type tmplData struct {
	CurrentCmd *Clic
	CalledCmds []*Clic
}

var tmplText = strings.TrimSpace(`
{{- $cmd := .CurrentCmd -}}
{{- $leftBrack := "[" -}}{{- $rightBrack := "]" -}}
Usage:

{{if .}}  {{end}}{{range $cmdIter := .CalledCmds}}
  {{- $cmdIter.FlagSet.Name}} {{if $cmdIter.FlagSet.Flags}}[FLAGS] {{end -}}
{{- end}}{{/* range .CalledCmd */}}
    {{- if $cmd.SubRequired}}{{$leftBrack = "{"}}{{$rightBrack = "}"}}{{end -}}
    {{- if $cmd.SubCmds}}{{$leftBrack}}{{end}}{{range $i, $subCmd := $cmd.SubCmds}}
      {{- if $subCmd.UsageConfig.Skip}}{{continue}}{{end}}
      {{- if $i}}|{{end}}{{$subCmd.FlagSet.Name}}
    {{- end}}{{/* range sub */}}
    {{- if and $cmd.SubCmds $cmd.OperandSet.Operands}}|{{end}}
    {{- range $i, $op := $cmd.OperandSet.Operands}}{{if $i}} {{end}}{{$op.Tag}}{{end}}
    {{- if $cmd.SubCmds}}{{$rightBrack}}{{end}}
    {{- if $cmd.UsageConfig.CmdDesc}}

      {{$cmd.UsageConfig.CmdDesc}}
    {{- end}}{{/* CmdDesc */}}
{{if $cmd.OperandSet.Operands}}
{{$cmd.OperandSet.Usage}}{{- end}}
{{- if $cmd.FlagSet.Flags}}
{{$cmd.FlagSet.Usage}}{{- end}}
`)

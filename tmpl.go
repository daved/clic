package clic

import (
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/daved/clic/tmpl"
)

// NewUsageTmpl returns the default TmplConfig value. This can be used
// as an example of how to setup custom usage output templating.
func NewUsageTmpl(c *Clic) *tmpl.Tmpl {
	type tmplData struct {
		ResolvedCmd    *Clic
		ResolvedCmdSet []*Clic
	}

	data := &tmplData{
		ResolvedCmd:    c,
		ResolvedCmdSet: resolvedCmdSet(c),
	}

	ensureSubCmdCatsSort(c)

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

		for _, sub := range cmd.SubCmds() {
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

	categoryLine := func(s string) string {
		name, desc, _ := strings.Cut(s, "|")
		return fmt.Sprintf("%-12s %s", name, desc)
	}

	subCmdsByCategoryFn := func(subs []*Clic, category string) []*Clic {
		return slices.DeleteFunc(slices.Clone(subs), func(c *Clic) bool {
			cat, _, _ := strings.Cut(category, "|")
			return c.HideUsage || c.Category != cat
		})
	}

	subCmdLine := func(c *Clic) string {
		return fmt.Sprintf("%-14s %s", c.FlagSet.Name(), c.Description)
	}

	fMap := template.FuncMap{
		"ResolvedCmdSetHint":  resolvedCmdSetHintFn,
		"SubsAndOperandsHint": subsAndOperandsHintFn,
		"StringsJoin":         strings.Join,
		"CategoryLine":        categoryLine,
		"SubCmdsByCategory":   subCmdsByCategoryFn,
		"SubCmdLine":          subCmdLine,
	}

	text := strings.TrimSpace(`
{{- $cmd := .ResolvedCmd -}}
Usage:

{{if .}}  {{end}}{{ResolvedCmdSetHint .ResolvedCmdSet}}{{SubsAndOperandsHint $cmd}}
    {{- if $cmd.Description}}

      {{$cmd.Description}}
    {{- end}}
{{if $cmd.FlagSet.Flags}}
{{$cmd.FlagSet.Usage}}{{- end}}
{{- if $cmd.Aliases}}
Aliases for {{$cmd.FlagSet.Name}}:

      {{StringsJoin $cmd.Aliases ", "}}
{{- end}}
{{- if $cmd.SubCmdCatsSort}}
Subcommands for {{$cmd.FlagSet.Name}}:

{{range $cmd.SubCmdCatsSort}}{{$catLine := CategoryLine .}}
{{- if $catLine}}  {{CategoryLine .}}{{- end}}
  {{- range $_, $sub := SubCmdsByCategory $cmd.SubCmds .}}
    {{SubCmdLine $sub}}
  {{- end}}

{{end}}
{{- end}}
`)

	return tmpl.New(text, fMap, data)
}

func ensureSubCmdCatsSort(c *Clic) {
	sort := slices.Clone(c.SubCmdCatsSort)
	for _, sub := range c.SubCmds() {
		if c.SubCmdCatsSort == nil && c.Category == "" {
			continue
		}

		if !slices.ContainsFunc(sort, func(s string) bool {
			prefix, _, _ := strings.Cut(s, "|")
			return prefix == sub.Category
		}) {
			sort = append(sort, sub.Category)
		}
	}

	if len(sort) > 0 {
		c.SubCmdCatsSort = sort
	}
}

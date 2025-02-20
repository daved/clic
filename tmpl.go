package clic

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/daved/flagset"
)

// Tmpl holds template configuration details.
type Tmpl struct {
	Text string
	FMap template.FuncMap
	Data any
}

func (t *Tmpl) Execute() (string, error) {
	tmpl := template.New("clic").Funcs(t.FMap)

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(t.Text)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(buf, t.Data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *Tmpl) String() string {
	s, err := t.Execute()
	if err != nil {
		s = fmt.Sprintf("%v\n", err)
	}
	return s
}

// NewUsageTmpl returns the default template configuration. This can be used as
// an example of how to setup custom usage output templating.
func NewUsageTmpl(c *Clic) *Tmpl {
	type tmplData struct {
		Cmd *Clic
	}

	data := &tmplData{
		Cmd: c,
	}

	cmdSetFn := func(c *Clic) []*Clic {
		all := []*Clic{c}

		for c.parent != nil {
			c = c.parent
			all = append(all, c)
		}

		slices.Reverse(all)

		return all
	}

	cmdSetHintFn := func(cmds []*Clic) string {
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
		out = pre + out + suf

		if out != "" {
			out = " " + out
		}
		return out
	}

	unhiddenFlagsFn := func(flags []*flagset.Flag) []*flagset.Flag {
		var out []*flagset.Flag
		for _, flag := range flags {
			if !flag.HideUsage {
				out = append(out, flag)
			}
		}
		return out
	}

	subCmdCatsSortFn := func(c *Clic) []string {
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
		return sort
	}

	categoryLine := func(s string) string {
		if s == "" {
			return ""
		}
		name, desc, _ := strings.Cut(s, "|")
		return fmt.Sprintf("%-18s %s", name, desc)
	}

	subCmdsByCategoryFn := func(subs []*Clic, category string) []*Clic {
		return slices.DeleteFunc(slices.Clone(subs), func(c *Clic) bool {
			cat, _, _ := strings.Cut(category, "|")
			return c.HideUsage || c.Category != cat
		})
	}

	subCmdLine := func(c *Clic) string {
		return fmt.Sprintf("%-18s %s", c.FlagSet.Name(), c.Description)
	}

	fMap := template.FuncMap{
		"CmdSet":              cmdSetFn,
		"CmdSetHint":          cmdSetHintFn,
		"SubsAndOperandsHint": subsAndOperandsHintFn,
		"UnhiddenFlags":       unhiddenFlagsFn,
		"StringsJoin":         strings.Join,
		"SubCmdCatsSort":      subCmdCatsSortFn,
		"CategoryLine":        categoryLine,
		"SubCmdsByCategory":   subCmdsByCategoryFn,
		"SubCmdLine":          subCmdLine,
	}

	text := strings.TrimSpace(`
{{- $cmd := .Cmd -}}
{{- $cmdSet := CmdSet $cmd -}}
{{- $unhiddenFlags := UnhiddenFlags $cmd.FlagSet.Flags -}}
{{- $subCmdCatsSort := SubCmdCatsSort $cmd -}}
{{if 1 -}}
Usage:

  {{CmdSetHint $cmdSet}}{{SubsAndOperandsHint $cmd}}
{{end -}}
{{if $cmd.Description}}
    {{$cmd.Description}}
{{end -}}
{{if $unhiddenFlags}}
{{$cmd.FlagSet.Usage -}}
{{end -}}
{{if $cmd.Aliases}}
Aliases for {{$cmd.FlagSet.Name}}:

      {{StringsJoin $cmd.Aliases ", "}}
{{end -}}
{{if $subCmdCatsSort}}
Subcommands for {{$cmd.FlagSet.Name}}:
{{range $subCmdCatsSort -}}{{- $catLine := CategoryLine . -}}
{{if $catLine}}
  {{$catLine}}{{end}}
{{range $_, $sub := SubCmdsByCategory $cmd.SubCmds . -}}
{{if 1}}{{end}}    {{SubCmdLine $sub}}
{{end -}}
{{if 1}}{{end -}}
{{end -}}
{{end -}}
`)

	return &Tmpl{text, fMap, data}
}

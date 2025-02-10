package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

// Tmpl tracks the template string and function map used for usage output
// templating.
type Tmpl struct {
	Text string
	FMap template.FuncMap
	Data any
}

func New(text string, fMap template.FuncMap, data any) *Tmpl {
	return &Tmpl{
		Text: text,
		FMap: fMap,
		Data: data,
	}
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

package gauge

import (
	"text/template"
	"bytes"
)

type TemplateInterface interface {
	Parse(c *Context) ([]byte, error)
}

type DefaultTemplate struct {
	template *template.Template
}

func newDefaultTemplate() *DefaultTemplate {
	tmpl :=  &DefaultTemplate{
		template: template.New("gauge"),
	}
	tmpl.SetTemplateLine(
		`{{if .HasMax}}{{.GetPercentage|printf "%3d"}}% {{end}}{{.GetBar}} {{.Status}} [{{.GetTime|printf "%.4f"}}sec|{{.GetMemory|printf "%.4f"}}MB] {{ .Extra }}`,
	)
	return tmpl
}

func (t *DefaultTemplate) SetTemplateLine(line string) {
	t.template = template.Must(t.template.Parse(line))
}

func (t *DefaultTemplate) Parse(c *Context) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := t.template.Execute(buff, c)
	if err != nil {
		return nil, err
	} else {
		return buff.Bytes(), nil
	}
}


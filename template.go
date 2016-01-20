package gauge

import (
	"text/template"
	"bytes"
)

type TemplateInterface interface {
	Parse(c *Context) ([]byte, error)
	AddFuncMap(funcs template.FuncMap)
}

type DefaultTemplate struct {
	template *template.Template
}

func newDefaultTemplate() *DefaultTemplate {
	tmpl :=  &DefaultTemplate{
		template: template.Must(template.New("gauge").Parse(
			` {{.Status}} {{.GetBar}}{{if .HasMax}} {{.GetPercentage|printf "%3d"}}%{{end}} [{{.GetTime.Seconds|printf "%05.2f"}}s{{if .HasMax}}/{{.GetEstimate|printf "%05.2f"}}s{{end}}][{{.GetMemory|printf "%05.2f"}}MB] {{ .Extra }}`,
		)),
	}
	return tmpl
}

func (t *DefaultTemplate) SetFormatLine(line string) {
	t.template = template.Must(t.template.Parse(line))
}

func (t *DefaultTemplate) AddFuncMap(funcs template.FuncMap) {
	t.template.Funcs(funcs)
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


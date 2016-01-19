package gauge

import (
	"io"
	"time"
	"fmt"
)

type Gauge struct {
	max 	 int
	cur 	 int
	writer 	 io.Writer
	start 	 time.Time
	template TemplateInterface
	context  *Context
}

// NewGauge creates a new progress gauge, if max
// is set 0 it will loop the bar until finished.
func NewGauge(max int, writer io.Writer) *Gauge {
	return &Gauge{
		max:    max,
		cur:    0,
		writer: writer,
	}
}

func (g *Gauge) Start() {
	g.start = time.Now()
	g.cur = 0
	g.write()
}

func (g *Gauge) SetTemplate(t TemplateInterface){
	g.template = t
}

func (g *Gauge) GetTemplate() TemplateInterface{
	if g.template == nil {
		g.template = newDefaultTemplate()
	}
	return g.template
}

func (g *Gauge) GetContext() *Context {
	if g.context == nil {
		g.context = &Context{gauge: g}
	}
	return g.context
}

func (g *Gauge) ClearLine() {
	fmt.Fprint(g.writer, "\033[2K\r")
}

func (g *Gauge) Advance(step int) {
	if g.max > 0 && g.cur + step > g.max {
		g.max = g.cur + step
	}
	g.cur += step
	g.ClearLine()
	g.write()
}

func (g *Gauge) Finished() {
	if g.max > g.cur {
		g.cur = g.max
	}

	g.ClearLine()
	g.write()

	if g.context.usage != nil {
		g.context.usage.Close()
	}

	fmt.Fprintf(g.writer, "\n")
}

func (g *Gauge) write() {
	buff, err := g.GetTemplate().Parse(g.GetContext())
	if err != nil {
		panic(err)
	}
	fmt.Fprint(g.writer, string(buff))
	g.GetContext().Extra = ""
}


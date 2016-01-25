package gauge

import (
	"io"
	"time"
	"fmt"
	"sync"
)

type Gauge struct {
	max 	 int
	cur 	 int
	writer 	 io.Writer
	start 	 time.Time
	template TemplateInterface
	context  *Context
	lock     sync.Mutex
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
	g.WriteLine()
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

func (g *Gauge) WriteLine() {
	g.lock.Lock()
	buff, err := g.GetTemplate().Parse(g.GetContext())
	if err != nil {
		panic(err)
	}
	g.write(string(buff))
	g.GetContext().Extra = ""
	g.lock.Unlock()
}


func (g *Gauge) ClearLine() {
	g.lock.Lock()
	g.write("\033[2K\r")
	g.lock.Unlock()
}

func (g *Gauge) Advance(step int) {
	if g.max > 0 && g.cur + step > g.max {
		g.max = g.cur + step
	}
	g.cur += step
	g.ClearLine()
	g.WriteLine()
}

func (g *Gauge) Finished() {
	if g.max > g.cur {
		g.cur = g.max
	}
	g.ClearLine()
	g.WriteLine()
	if g.context.usage != nil {
		g.context.usage.Close()
	}
	g.write("\n")
}

func (g *Gauge) write(line string) {
	fmt.Fprint(g.writer, line)
}


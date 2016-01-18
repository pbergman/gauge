package gauge

import (
	"time"
	"fmt"
	"strings"
	"strconv"
	"github.com/pbergman/gauge/debug"
)

type Context struct {
	Extra string
	gauge *Gauge
	usage *debug.MemoryUsage
}

func (c *Context) HasMax() bool {
	return c.gauge.max > 0
}

func (c *Context) GetTime() float64 {
	return time.Now().Sub(c.gauge.start).Seconds()
}

func (c *Context) GetPercentage() int {
	return int(float64(c.gauge.cur) * (100 / float64(c.gauge.max)))
}

func (c *Context) GetMemory() float64 {
	if c.usage == nil {
		c.usage, _ = debug.NewSelfMemoryUsage()
	}
	size, _ := c.usage.GetSize()
	return float64(size)/1024
}

func (c *Context) Status() string {
	if c.gauge.max <= 0 {
		return fmt.Sprintf("%d", c.gauge.cur)
	} else {
		str_max := strconv.Itoa(c.gauge.max)
		str_now :=  strconv.Itoa(c.gauge.cur)
		prefix := strings.Repeat(" ", len(str_max) - len(str_now))
		return fmt.Sprintf("%s%s/%s", prefix, str_now, str_max)
	}
}

func (c *Context) GetBar() string {
	var width float64 = 28
	line := strings.Repeat("-", int(width))
	if c.gauge.max <= 0 {
		pos := c.gauge.cur % int(width)
		if pos == 0 {
			return "[" + line + "]"
		} else {
			return "[" + line[:pos - 1] + ">" + line[pos:] + "]"
		}
	} else {
		pos  := (width/float64(c.gauge.max)) * float64(c.gauge.cur)
		if pos > width {
			pos = width
		}
		if pos - 1 <= 0 {
			return "[" + line + "]"
		} else {
			return "[" + line[:int(pos) - 1] + ">" + line[int(pos):] + "]"
		}
	}
	return ""
}
##GOUGE 

A simple progress gauge for displaying progress status.

When the first argument of the NewGauge method is set 0, the progress bar will loop until finished.

example:

```
func main(){
	var wg sync.WaitGroup
	gauge := gauge.NewGauge(5, os.Stdout)
	gauge.Start()
	
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			// some work .....
			//
			// if you want to add extra context 
			// gauge.GetContext().Extra = "finished ....."
			gauge.Advance(1)
			wg.Done()
		}()
	}
	
	wg.Wait()
	gauge.Finished()
	
	// When finished should print something like: 
	// 100% [--------------------------->] 6/6 [6.0019sec|3.0352MB] 
```


The output line can be overwritten by implementing the TemplateInterface:

```
type templ struct {
	template *template.Template
}

func newTempl() *templ{
	return &templ{
		template: template.Must(template.New("gauge").Parse(
			`{{if .HasMax}}{{.GetPercentage|printf "%3d"}}% {{end}}{{.GetBar}} {{.Status}}`,
		)),
	}
}

func (t *templ) Parse(c *gauge.Context) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := t.template.Execute(buff, c)
	if err != nil {
		return nil, err
	} else {
		return buff.Bytes(), nil
	}
}

func main(){
	gauge := gauge.NewGauge(5, os.Stdout)
	gauge.SetTemplate(newTempl())
	gauge.Start()
.....
	
```
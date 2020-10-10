package sendo

import (
	"github.com/kamva/tracer"
	"io/ioutil"
	"text/template"
)

// parseTextTemplates parses list of text templates
func parseTextTemplates(rootName string, l map[string]string) (*template.Template, error) {
	t := template.New(rootName)

	for name, path := range l {
		c, err := fileContent(path)
		if err != nil {
			return nil, tracer.Trace(err)
		}
		if _, err := t.New(name).Parse(c); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	return t, nil
}

func fileContent(fileName string) (string, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(b), tracer.Trace(err)
}

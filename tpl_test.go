package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlTemplate(t *testing.T) {
	type io struct {
		Input    string
		Template string
		Output   string
	}

	tests := []io{
		io{
			Input:    "test: value",
			Template: "{{.test}}",
			Output:   "value",
		},
		io{
			Input:    "name: Max\nage: 15",
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
		},
		io{
			Input:    "legumes:\n  - potato\n  - onion\n  - cabbage",
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
		},
	}

	for _, test := range tests {
		tpl_file, err := ioutil.TempFile("", "")
		assert.Nil(t, err)
		defer func() { os.Remove(tpl_file.Name()) }()

		_, err = tpl_file.WriteString(test.Template)
		assert.Nil(t, err)
		tpl_file.Close()

		output := bytes.NewBuffer(nil)
		err = ExecuteTemplates(
			[]string{},
			strings.NewReader(test.Input),
			output,
			tpl_file.Name(),
		)
		assert.Nil(t, err)

		assert.Equal(t, test.Output, output.String())

	}
}

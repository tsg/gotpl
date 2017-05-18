package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	type io struct {
		Input    string
		Template string
		Output   string
		Format   string
	}

	tests := []io{
		io{
			Input:    "test: value",
			Template: "{{.test}}",
			Output:   "value",
			Format:   "yaml",
		},
		io{
			Input:    "name: Max\nage: 15",
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
			Format:   "yaml",
		},
		io{
			Input:    "legumes:\n  - potato\n  - onion\n  - cabbage",
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
			Format:   "yaml",
		},
		io{
			Input:    "{\"test\": \"value\"}",
			Template: "{{.test}}",
			Output:   "value",
			Format:   "json",
		},
		io{
			Input:    "{\"name\": \"Max\", \"age\": 15}",
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
			Format:   "json",
		},
		io{
			Input:    "{\"legumes\": [\"potato\", \"onion\", \"cabbage\"]}",
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
			Format:   "json",
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
		values, err := ParseValues(strings.NewReader(test.Input), test.Format)
		assert.Nil(t, err)
		err = ExecuteTemplates(values, output,
			tpl_file.Name())
		assert.Nil(t, err)

		assert.Equal(t, test.Output, output.String())

	}
}

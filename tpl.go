package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values_in io.Reader, out io.Writer, tpl_files string) error {
	var r io.Reader
	r, err := os.Open(tpl_files)
	if err != nil {
		return err
	}
	templateData, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Could not read input: %v", err)
	}

	tpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, values_in)
	if err != nil {
		return fmt.Errorf("Failed to read standard input: %v", err)
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}

	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func main() {
	err := ExecuteTemplates(os.Stdin, os.Stdout, os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values map[string]interface{}, out io.Writer, tpl_files ...string) error {
	tpl, err := template.ParseFiles(tpl_files...)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func ParseValues(values_in io.Reader, format string) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, values_in)
	if err != nil {
		return nil, fmt.Errorf("Failed to read standard input: %v", err)
	}

	var values map[string]interface{}

	switch format {
	case "json":
		err = json.Unmarshal(buf.Bytes(), &values)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse standard input: %v", err)
		}
	case "yaml":
		err = yaml.Unmarshal(buf.Bytes(), &values)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse standard input: %v", err)
		}
	default:
		return nil, fmt.Errorf("Unknown format: %s", format)
	}

	return values, nil

}

func main() {
	format := flag.String("format", "yaml", "format of input values: yaml or json")
	flag.Parse()

	values, err := ParseValues(os.Stdin, *format)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = ExecuteTemplates(values, os.Stdout, flag.Args()...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

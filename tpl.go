package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values_in io.Reader, out io.Writer, tpl_files ...string) error {
	tpl, err := template.ParseFiles(tpl_files...)
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

	values["Env"] = getEnvironment(os.Environ())

	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}

func getEnvironment(data []string) map[string]string {
	items := make(map[string]string)
	for _, item := range data {
		key, val := getKeyVal(item)
		items[key] = val
	}
	return items
}

func getKeyVal(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = splits[1]
	return
}

func main() {
	err := ExecuteTemplates(os.Stdin, os.Stdout, os.Args[1:]...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

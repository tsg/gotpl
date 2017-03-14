package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	docopt "github.com/docopt/docopt-go"

	"gopkg.in/yaml.v2"
)

const Version = "0.0.1"

// Usage is a docopt-formatted specification for this application's command line interface.
const Usage = `Usage:
  gotpl [-e=<kv>...] <template>...
  gotpl -h
  gotpl -v

Options:
  -e=<kv>       Supply a YAML key value pair, like "ticks: 3"
  -h --help     Show usage information
  -v --version  Show version information`

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(flag_values_in []string, values_in io.Reader, out io.Writer, tpl_files ...string) error {
	tpl, err := template.ParseFiles(tpl_files...)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	buf := bytes.NewBuffer(nil)

	for _, flag_value := range flag_values_in {
		buf.Write([]byte(flag_value))
	}

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
	arguments, _ := docopt.Parse(Usage, nil, true, Version, false)

	var flag_values []string

	if arguments["-e"] != nil {
		flag_values = arguments["-e"].([]string)
	}

	err := ExecuteTemplates(flag_values, os.Stdin, os.Stdout, arguments["<template>"].([]string)...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

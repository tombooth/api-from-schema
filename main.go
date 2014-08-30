package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"

	"github.com/docopt/docopt-go"
	"github.com/tombooth/api-from-schema/schematic"
)

func main() {
	usage := `Api from schema

Usage:
  api-from-schema <json_schema>
  api-from-schema -h | --help
  api-from-schema --version

Options:
  -h --help     Show this screen.
  --version     Show version.
`

	arguments, _ := docopt.Parse(usage, nil, true, "Api from schema 0.1.0", false)
	path := arguments["<json_schema>"].(string)

	if apiSchema, err := schema.ParseSchema(path); err == nil {
		apiSchema.Resolve(nil)
		models := ModelsFromSchema(apiSchema)

		context := struct {
			Models []Model
		}{
			Models: models,
		}

		var apiSource bytes.Buffer

		apiTmpl, _ := template.ParseFiles("templates/api.tmpl")
		apiTmpl.Execute(&apiSource, context)

		if formattedSource, err := format.Source(apiSource.Bytes()); err != nil {
			fmt.Println("Failed to format source: %v", err)
		} else {
			os.Stdout.Write(formattedSource)
		}
	}
}

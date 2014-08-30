package main

import (
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
		endpoints := EndpointsFromSchema(apiSchema)

		apiTmpl, _ := template.ParseFiles("templates/api.tmpl")
		apiTmpl.Execute(os.Stdout, endpoints)
	}
}

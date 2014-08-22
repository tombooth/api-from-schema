package main

import (
	"fmt"
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
	if schema_rep, err := schema.ParseSchema(path); err == nil {
		fmt.Println("%s", schema_rep)
	}
}

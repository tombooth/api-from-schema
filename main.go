package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/tombooth/api-from-schema/schematic"
)

func main() {
	usage := `Api from schema

Usage:
  api-from-schema [options] <json_schema>
  api-from-schema -h | --help
  api-from-schema --version

Options:
  -h --help               Show this screen.
  --version               Show version.
  --templates=<path>      Path to template directory [default: templates/]
`

	arguments, _ := docopt.Parse(usage, nil, true, "Api from schema 0.1.0", false)
	path := arguments["<json_schema>"].(string)
	templateStore, err := NewTemplateStore(arguments["--templates"].(string))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load templates: %v", err)
		return
	}

	if apiSchema, err := schema.ParseSchema(path); err == nil {
		apiSchema.Resolve(nil)
		models := ModelsFromSchema(apiSchema)

		context := struct {
			Models        []Model
			TemplateStore TemplateStore
		}{
			Models:        models,
			TemplateStore: templateStore,
		}

		if apiOutput, err := templateStore.ExecuteAndFormat(context, "api.tmpl", "handlerfunc.tmpl"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute api template: %v", err)
		} else {
			fmt.Fprint(os.Stdout, apiOutput)
		}
	}
}

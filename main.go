package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/tombooth/api-from-schema/schematic"
)

type Endpoint struct {
	URL                string
	Method             string
	IsList             bool
	ResponseDefinition *schema.Schema
}

func EndpointsFromSchema(apiSchema *schema.Schema) []Endpoint {
	endpoints := []Endpoint{}

	for _, definition := range apiSchema.Definitions {
		for _, link := range definition.Links {
			endpoint := Endpoint{
				URL:                link.HRef.URLPattern(),
				Method:             link.Method,
				IsList:             link.Rel == "instances",
				ResponseDefinition: definition,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

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

		for _, endpoint := range EndpointsFromSchema(apiSchema) {
			fmt.Println(endpoint.Method, endpoint.URL, endpoint.IsList)
		}
	}
}

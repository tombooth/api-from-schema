package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/tombooth/api-from-schema/schematic"
)

type Endpoint struct {
	URL    string
	Method string
	IsList bool

	ResponseDefinition *schema.Schema
	HRefDefinition     *schema.HRef
}

func EndpointsFromSchema(apiSchema *schema.Schema) ([]Endpoint, map[*schema.Schema][]Endpoint) {
	endpoints := []Endpoint{}
	definitionToEndpoints := make(map[*schema.Schema][]Endpoint)

	for _, definition := range apiSchema.Definitions {
		endpointsForDefinition := []Endpoint{}

		for _, link := range definition.Links {
			endpoint := Endpoint{
				URL:    link.HRef.URLPattern(),
				Method: link.Method,
				IsList: link.Rel == "instances",

				ResponseDefinition: definition,
				HRefDefinition:     link.HRef,
			}
			endpointsForDefinition = append(endpointsForDefinition, endpoint)
		}

		endpoints = append(endpoints, endpointsForDefinition...)
		definitionToEndpoints[definition] = endpointsForDefinition
	}

	return endpoints, definitionToEndpoints
}

func (endpoint *Endpoint) HandlerName() string {
	multiple := ""
	if endpoint.IsList {
		multiple = "s"
	}

	parameters := ""
	if len(endpoint.HRefDefinition.Order) > 0 {
		parameters = "By"
		for _, name := range endpoint.HRefDefinition.Order {
			parameters += strings.Title(name)
		}
	}

	return strings.Join([]string{
		endpoint.Method,
		strings.Title(endpoint.ResponseDefinition.Title),
		multiple,
		parameters,
	}, "")
}

func (endpoint *Endpoint) HandlerDefinition() string {
	var handlerDefinition bytes.Buffer

	apiTmpl, _ := template.ParseFiles("templates/handlerfunc.tmpl")
	apiTmpl.Execute(&handlerDefinition, endpoint)

	return handlerDefinition.String()
}

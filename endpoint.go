package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tombooth/api-from-schema/schematic"
)

type Endpoint struct {
	URL    string
	Method string
	IsList bool

	ResponseDefinition *schema.Schema
	HRefDefinition     *schema.HRef
}

type EndpointSignature struct {
	IsList bool
	Vars   string
}

func EndpointFromLink(link schema.Link, parent *schema.Schema) Endpoint {
	return Endpoint{
		URL:    link.HRef.URLPattern(),
		Method: link.Method,
		IsList: link.Rel == "instances",

		ResponseDefinition: parent,
		HRefDefinition:     link.HRef,
	}
}

func EndpointsFromDefinition(definition *schema.Schema) []Endpoint {
	endpointsForDefinition := []Endpoint{}

	for _, link := range definition.Links {
		endpointsForDefinition = append(
			endpointsForDefinition,
			EndpointFromLink(link, definition))
	}

	return endpointsForDefinition
}

func EndpointsFromSchema(apiSchema *schema.Schema) []Endpoint {
	endpoints := []Endpoint{}

	for _, definition := range apiSchema.Definitions {
		endpoints = append(
			endpoints,
			EndpointsFromDefinition(definition)...)
	}

	return endpoints
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

func (endpoint *Endpoint) HandlerDefinition(model Model, templateStore TemplateStore) string {
	handlerDefinition, err := templateStore.Execute(struct {
		Model    *Model
		Endpoint *Endpoint
	}{
		Model:    &model,
		Endpoint: endpoint,
	}, "handlerfunc.tmpl")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute handler template: %v", err)
	}

	return handlerDefinition
}

func (endpoint *Endpoint) Vars() []string {
	return endpoint.HRefDefinition.Order
}

func (endpoint *Endpoint) RequiresModel() bool {
	return len(endpoint.Vars()) > 0 || endpoint.IsList
}

func (endpoint *Endpoint) Signature() EndpointSignature {
	return EndpointSignature{
		IsList: endpoint.IsList,
		Vars:   strings.Join(endpoint.Vars(), ""),
	}
}

package main

import (
	"strings"

	"github.com/tombooth/api-from-schema/schematic"
)

type Endpoint struct {
	URL    string
	Method string
	IsList bool

	Model Model

	hrefDefinition *schema.HRef
}

type EndpointSignature struct {
	IsList bool
	Vars   string
}

func EndpointFromLink(link schema.Link, parent Model) Endpoint {
	return Endpoint{
		URL:    link.HRef.URLPattern(),
		Method: link.Method,
		IsList: link.Rel == "instances",

		Model: parent,

		hrefDefinition: link.HRef,
	}
}

func EndpointsFromModel(model Model) []Endpoint {
	endpointsForModel := []Endpoint{}

	for _, link := range model.Definition.Links {
		endpointsForModel = append(
			endpointsForModel,
			EndpointFromLink(link, model))
	}

	return endpointsForModel
}

func (endpoint *Endpoint) HandlerName() string {
	multiple := ""
	if endpoint.IsList {
		multiple = "s"
	}

	parameters := ""
	if len(endpoint.hrefDefinition.Order) > 0 {
		parameters = "By"
		for _, name := range endpoint.hrefDefinition.Order {
			parameters += strings.Title(name)
		}
	}

	return strings.Join([]string{
		endpoint.Method,
		strings.Title(endpoint.Model.Definition.Title),
		multiple,
		parameters,
	}, "")
}

func (endpoint *Endpoint) Vars() []string {
	return endpoint.hrefDefinition.Order
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

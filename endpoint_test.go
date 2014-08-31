package main

import (
	"testing"

	"github.com/tombooth/api-from-schema/schematic"
)

func assertEndpoint(t *testing.T, endpoint Endpoint, expectedURL string, expectedMethod string, isList bool) {
	if endpoint.URL != expectedURL {
		t.Errorf("URL should be %v, got %v", expectedURL, endpoint.URL)
	}
	if endpoint.Method != expectedMethod {
		t.Errorf("Method should be %v, got %v", expectedMethod, endpoint.Method)
	}
	if endpoint.IsList != isList {
		t.Errorf("IsList should be %v, got %v", isList, endpoint.IsList)
	}
}

func endpointsFromJSON(path string) []Endpoint {
	apiSchema, _ := schema.ParseSchema("fixtures/test-api.json")
	apiSchema.Resolve(nil)

	models := ModelsFromSchema(apiSchema)
	endpoints := []Endpoint{}
	for _, model := range models {
		endpoints = append(endpoints, model.Endpoints...)
	}

	return endpoints
}

func TestEndpoints(t *testing.T) {
	endpoints := endpointsFromJSON("fixtures/test-api.json")

	assertEndpoint(t, endpoints[0], "/thing/{thingIdentifier:[0-9]+}", "GET", false)
	assertEndpoint(t, endpoints[1], "/thing", "GET", true)
}

func TestHandlerName(t *testing.T) {
	endpoints := endpointsFromJSON("fixtures/test-api.json")

	if endpoints[0].HandlerName() != "GETThingByThingIdentifier" ||
		endpoints[1].HandlerName() != "GETThings" {
		t.Errorf("HanderName failed. [0] %v [1] %v", endpoints[0].HandlerName(), endpoints[1].HandlerName())
	}
}

func expectRequiresModel(t *testing.T, expected bool, order []string, isList bool) {
	endpoint := Endpoint{
		hrefDefinition: &schema.HRef{
			Order: order,
		},
		IsList: isList,
	}
	assertEqual(t, expected, endpoint.RequiresModel())
}

func TestRequiresModel(t *testing.T) {
	expectRequiresModel(t, true, []string{"foo"}, false)
	expectRequiresModel(t, true, []string{"foo"}, true)
	expectRequiresModel(t, true, []string{}, true)
	expectRequiresModel(t, false, []string{}, false)
}

func TestSignature(t *testing.T) {
	endpoint := Endpoint{
		hrefDefinition: &schema.HRef{
			Order: []string{"foo", "bar"},
		},
		IsList: true,
	}
	assertEqual(t, "foobar", endpoint.Signature().Vars)
	assertEqual(t, true, endpoint.Signature().IsList)
}

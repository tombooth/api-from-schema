package main

import (
	"github.com/tombooth/api-from-schema/schematic"
	"testing"
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

func TestEndpoints(t *testing.T) {
	apiSchema, _ := schema.ParseSchema("fixtures/test-api.json")
	apiSchema.Resolve(nil)

	endpoints := EndpointsFromSchema(apiSchema)

	assertEndpoint(t, endpoints[0], "/thing/{thingIdentifier:[0-9]+}", "GET", false)
	assertEndpoint(t, endpoints[1], "/thing", "GET", true)
}

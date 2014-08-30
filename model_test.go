package main

import (
	"testing"

	"github.com/tombooth/api-from-schema/schematic"
)

func modelsFromJSON(path string) []Model {
	apiSchema, _ := schema.ParseSchema(path)
	apiSchema.Resolve(nil)

	return ModelsFromSchema(apiSchema)
}

func modelByName(models []Model, name string) *Model {
	for _, model := range models {
		if model.Name == name {
			return &model
		}
	}
	return nil
}

func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	if got != expected {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func expectConstructor(t *testing.T, constructors []Constructor, name, arguments, returnType string) {
	for _, constructor := range constructors {
		if constructor.Name == name {
			assertEqual(t, arguments, constructor.Arguments)
			assertEqual(t, returnType, constructor.ReturnType)
			return
		}
	}
	t.Errorf("Failed to find a constructor called %v, in %v", name, constructors)
}

func TestConstructor(t *testing.T) {
	models := modelsFromJSON("fixtures/model-tests.json")

	assertEqual(t, 3, len(models))

	constructors := modelByName(models, "Thing").Constructors()
	assertEqual(t, 3, len(constructors))
	expectConstructor(
		t, constructors,
		"ThingByThingIdentifier",
		"thingIdentifier string",
		"(Thing, error)")
	expectConstructor(
		t, constructors,
		"ListThingsByThing1IdentifierAndThing2Identifier",
		"thing1Identifier, thing2Identifier string",
		"([]Thing, error)")
	expectConstructor(t, constructors,
		"ListThings",
		"",
		"([]Thing, error)")

	constructors = modelByName(models, "Thing1").Constructors()
	assertEqual(t, 0, len(constructors))

	constructors = modelByName(models, "Thing2").Constructors()
	assertEqual(t, 1, len(constructors))
	expectConstructor(t, constructors,
		"Thing2ByThing2Identifier",
		"thing2Identifier string",
		"(Thing2, error)")
}

func TestConstructorForEndpoint(t *testing.T) {
	models := modelsFromJSON("fixtures/model-tests.json")

	assertEqual(t, 3, len(models))

	model := modelByName(models, "Thing")
	_, err := model.ConstructorForEndpoint(model.Endpoints[0])
	assertEqual(t, nil, err)

	model = modelByName(models, "Thing1")
	_, err = model.ConstructorForEndpoint(model.Endpoints[0])
	assertEqual(t, true, err != nil)
}

package main

import (
	"github.com/tombooth/api-from-schema/schematic"
)

type Model struct {
	Name       string
	Definition *schema.Schema

	Endpoints []Endpoint
}

func ModelsFromSchema(apiSchema *schema.Schema) []Model {
	models := []Model{}

	for _, definition := range apiSchema.Definitions {
		model := Model{
			Name:       definition.Title,
			Definition: definition,

			Endpoints: EndpointsFromDefinition(definition),
		}
		models = append(models, model)
	}

	return models
}

func (model *Model) AsType() string {
	return model.Definition.GoType()
}

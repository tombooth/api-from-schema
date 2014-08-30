package main

import (
	"github.com/tombooth/api-from-schema/schematic"
)

type Model struct {
	Name       string
	Definition *schema.Schema

	Endpoints []Endpoint
}

func ModelsFromDTE(DTE map[*schema.Schema][]Endpoint) []Model {
	models := []Model{}

	for definition, endpoints := range DTE {
		model := Model{
			Name:       definition.Title,
			Definition: definition,

			Endpoints: endpoints,
		}
		models = append(models, model)
	}

	return models
}

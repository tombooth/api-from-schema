package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tombooth/api-from-schema/schematic"
)

type Model struct {
	Name       string
	Definition *schema.Schema

	Endpoints []Endpoint
}

type Constructor struct {
	Name        string
	Arguments   []string
	ReturnType  string
	ReturnsList bool
}

func ModelsFromSchema(apiSchema *schema.Schema) []Model {
	models := []Model{}

	for _, definition := range apiSchema.Definitions {
		model := Model{
			Name:       definition.Title,
			Definition: definition,
		}
		model.Endpoints = EndpointsFromModel(model)
		models = append(models, model)
	}

	return models
}

func (model *Model) AsType() string {
	return model.Definition.GoType()
}

func (model *Model) ConstructorForEndpoint(endpoint Endpoint) (Constructor, error) {

	if !endpoint.RequiresModel() {
		return Constructor{}, errors.New("Cannot create a constructor object for an endpoint that doesn't need one")
	}

	vars := endpoint.Vars()

	namePrefix := ""
	returnTypePrefix := ""
	if endpoint.IsList {
		namePrefix = "List"
		returnTypePrefix = "[]"
	}

	nameSuffix := ""
	if endpoint.IsList {
		nameSuffix = nameSuffix + "s"
	}
	if len(vars) > 0 {
		nameSuffix = nameSuffix + "By"
		byThings := []string{}
		for _, v := range vars {
			byThings = append(byThings, strings.Title(v))
		}
		nameSuffix = nameSuffix + strings.Join(byThings, "And")
	}

	return Constructor{
		Name:        fmt.Sprintf("%v%v%v", namePrefix, model.Name, nameSuffix),
		Arguments:   vars,
		ReturnType:  fmt.Sprintf("(%v%v, error)", returnTypePrefix, model.Name),
		ReturnsList: endpoint.IsList,
	}, nil
}

func (model *Model) Constructors() []Constructor {
	constructors := []Constructor{}
	signatures := make(map[EndpointSignature]bool)

	for _, endpoint := range model.Endpoints {
		if endpoint.RequiresModel() {
			signature := endpoint.Signature()
			if !signatures[signature] {
				// it shouldn't through an error as we have already validated it will have
				// a constructor
				constructor, _ := model.ConstructorForEndpoint(endpoint)
				constructors = append(constructors, constructor)
				signatures[signature] = true
			}
		}
	}

	return constructors
}

func (constructor *Constructor) ArgumentsAsString() string {
	arguments := ""
	if len(constructor.Arguments) > 0 {
		arguments = strings.Join(constructor.Arguments, ", ") + " string"
	}
	return arguments
}

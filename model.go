package main

import (
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
	Arguments   string
	ReturnType  string
	ReturnsList bool
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

func (model *Model) endpointsWithRetrieval() []Endpoint {
	endpoints := []Endpoint{}
	for _, endpoint := range model.Endpoints {
		if len(endpoint.Vars()) > 0 || endpoint.IsList {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func (model *Model) Constructors() []Constructor {
	constructors := []Constructor{}

	type constructorSignature struct {
		IsList bool
		Vars   string
	}

	signatures := make(map[constructorSignature][]string)

	for _, endpoint := range model.endpointsWithRetrieval() {
		signature := constructorSignature{
			IsList: endpoint.IsList,
			Vars:   strings.Join(endpoint.Vars(), ""),
		}
		if _, found := signatures[signature]; !found {
			signatures[signature] = endpoint.Vars()
		}
	}

	for signature, vars := range signatures {
		namePrefix := ""
		returnTypePrefix := ""
		if signature.IsList {
			namePrefix = "List"
			returnTypePrefix = "[]"
		}

		nameSuffix := ""
		arguments := ""
		if signature.IsList {
			nameSuffix = nameSuffix + "s"
		}
		if len(vars) > 0 {
			nameSuffix = nameSuffix + "By"
			byThings := []string{}
			for _, v := range vars {
				byThings = append(byThings, strings.Title(v))
			}
			nameSuffix = nameSuffix + strings.Join(byThings, "And")
			arguments = strings.Join(vars, ", ") + " string"
		}

		constructor := Constructor{
			Name:        fmt.Sprintf("%v%v%v", namePrefix, model.Name, nameSuffix),
			Arguments:   arguments,
			ReturnType:  fmt.Sprintf("(%v%v, error)", returnTypePrefix, model.Name),
			ReturnsList: signature.IsList,
		}
		constructors = append(constructors, constructor)
	}

	return constructors
}

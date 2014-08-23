package schema

import (
	"reflect"
	"strings"
	"testing"
)

var typeTests = []struct {
	Schema *Schema
	Type   string
}{
	{
		Schema: &Schema{
			Type: "boolean",
		},
		Type: "bool",
	},
	{
		Schema: &Schema{
			Type: "number",
		},
		Type: "float64",
	},
	{
		Schema: &Schema{
			Type: "integer",
		},
		Type: "int",
	},
	{
		Schema: &Schema{
			Type: "any",
		},
		Type: "interface{}",
	},
	{
		Schema: &Schema{
			Type: "string",
		},
		Type: "string",
	},
	{
		Schema: &Schema{
			Type:   "string",
			Format: "date-time",
		},
		Type: "time.Time",
	},
	{
		Schema: &Schema{
			Type: []interface{}{"null", "string"},
		},
		Type: "*string",
	},
	{
		Schema: &Schema{
			Type: "array",
			Items: &Schema{
				Type: "string",
			},
		},
		Type: "[]string",
	},
	{
		Schema: &Schema{
			Type:   []interface{}{"null", "string"},
			Format: "date-time",
		},
		Type: "*time.Time",
	},
}

func TestSchemaType(t *testing.T) {
	for i, tt := range typeTests {
		kind := tt.Schema.GoType()
		if !strings.Contains(kind, tt.Type) {
			t.Errorf("%d: wants %v, got %v", i, tt.Type, kind)
		}
	}
}

var regexTests = []struct {
	Schema *Schema
	Regex  string
}{
	{
		Schema: &Schema{
			Type: "boolean",
		},
		Regex: "true|false",
	},
	{
		Schema: &Schema{
			Type: "number",
		},
		Regex: "[0-9\\.]+",
	},
	{
		Schema: &Schema{
			Type: "integer",
		},
		Regex: "[0-9]+",
	},
	{
		Schema: &Schema{
			Type: "any",
		},
		Regex: ".*",
	},
	{
		Schema: &Schema{
			Type: "string",
		},
		Regex: ".*",
	},
	{
		Schema: &Schema{
			Type: []interface{}{"string", "integer"},
		},
		Regex: "(.*)|([0-9]+)",
	},
	{
		Schema: &Schema{
			Type:    "string",
			Pattern: "[a-z]+",
		},
		Regex: "[a-z]+",
	},
	{
		Schema: &Schema{
			Type:   "string",
			Format: "uuid",
		},
		Regex: "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	},
	{
		Schema: &Schema{
			Type:   "string",
			Format: "date-time",
		},
		Regex: "([\\+-]?\\d{4}(?!\\d{2}\\b))((-?)((0[1-9]|1[0-2])(\\3([12]\\d|0[1-9]|3[01]))?|W([0-4]\\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\\d|[12]\\d{2}|3([0-5]\\d|6[1-6])))([T\\s]((([01]\\d|2[0-3])((:?)[0-5]\\d)?|24\\:?00)([\\.,]\\d+(?!:))?)?(\\17[0-5]\\d([\\.,]\\d+)?)?([zZ]|([\\+-])([01]\\d|2[0-3]):?([0-5]\\d)?)?)?)?",
	},
}

func TestSchemaRegex(t *testing.T) {
	for i, tt := range regexTests {
		kind := tt.Schema.Regex()
		if !strings.Contains(kind, tt.Regex) {
			t.Errorf("%d: wants %v, got %v", i, tt.Regex, kind)
		}
	}
}

var linkTests = []struct {
	Link *Link
	Type string
}{}

func TestLinkType(t *testing.T) {
	for i, lt := range linkTests {
		kind := lt.Link.GoType()
		if !strings.Contains(kind, lt.Type) {
			t.Errorf("%d: wants %v, got %v", i, lt.Type, kind)
		}
	}
}

var paramsTests = []struct {
	Schema     *Schema
	Link       *Link
	Order      []string
	Parameters map[string]string
}{
	{
		Schema: &Schema{},
		Link: &Link{
			HRef: NewHRef("/destroy/"),
			Rel:  "destroy",
		},
		Parameters: map[string]string{},
	},
	{
		Schema: &Schema{},
		Link: &Link{
			HRef: NewHRef("/instances/"),
			Rel:  "instances",
		},
		Order:      []string{"lr"},
		Parameters: map[string]string{"lr": "*ListRange"},
	},
	{
		Schema: &Schema{},
		Link: &Link{
			Rel:  "update",
			HRef: NewHRef("/update/"),
			Schema: &Schema{
				Type: "string",
			},
		},
		Order:      []string{"o"},
		Parameters: map[string]string{"o": "string"},
	},
	{
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"struct": {
					Definitions: map[string]*Schema{
						"uuid": {
							Type: "string",
						},
					},
				},
			},
		},
		Link: &Link{
			HRef: NewHRef("/results/{(%23%2Fdefinitions%2Fstruct%2Fdefinitions%2Fuuid)}"),
		},
		Order:      []string{"structUUID"},
		Parameters: map[string]string{"structUUID": "string"},
	},
}

func TestParameters(t *testing.T) {
	for i, pt := range paramsTests {
		pt.Link.Resolve(pt.Schema)
		order, params := pt.Link.Parameters()
		if !reflect.DeepEqual(order, pt.Order) {
			t.Errorf("%d: wants %v, got %v", i, pt.Order, order)
		}
		if !reflect.DeepEqual(params, pt.Parameters) {
			t.Errorf("%d: wants %v, got %v", i, pt.Parameters, params)
		}

	}
}

var valuesTests = []struct {
	Schema *Schema
	Name   string
	Link   *Link
	Values []string
}{
	{
		Schema: &Schema{},
		Name:   "Result",
		Link: &Link{
			Rel: "destroy",
		},
		Values: []string{"error"},
	},
	{
		Schema: &Schema{},
		Name:   "Result",
		Link: &Link{
			Rel: "instances",
		},
		Values: []string{"[]*Result", "error"},
	},
	{
		Schema: &Schema{
			Type: "object",
			Properties: map[string]*Schema{
				"value": {
					Type: "integer",
				},
			},
			Required: []string{"value"},
		},
		Name: "Result",
		Link: &Link{
			Rel: "self",
		},
		Values: []string{"*Result", "error"},
	},
	{
		Schema: &Schema{
			Type:                 "object",
			AdditionalProperties: false,
			PatternProperties: map[string]*Schema{
				"^\\w+$": {
					Type: "string",
				},
			},
		},
		Name: "ConfigVar",
		Link: &Link{
			Rel: "self",
		},
		Values: []string{"map[string]string", "error"},
	},
}

func TestValues(t *testing.T) {
	for i, vt := range valuesTests {
		values := vt.Schema.Values(vt.Name, vt.Link)
		if !reflect.DeepEqual(values, vt.Values) {
			t.Errorf("%d: wants %v, got %v", i, vt.Values, values)
		}
	}
}

package schema

import (
	"reflect"
	"testing"
)

var resolveTests = []struct {
	Ref      string
	Schema   *Schema
	Resolved *Schema
}{
	{
		Ref: "#/definitions/uuid",
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"uuid": {
					Title: "Identifier",
				},
			},
		},
		Resolved: &Schema{
			Title: "Identifier",
		},
	},
	{
		Ref: "#/definitions/struct/definitions/uuid",
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"struct": {
					Definitions: map[string]*Schema{
						"uuid": {
							Title: "Identifier",
						},
					},
				},
			},
		},
		Resolved: &Schema{
			Title: "Identifier",
		},
	},
}

func TestReferenceResolve(t *testing.T) {
	for i, rt := range resolveTests {
		ref := Reference(rt.Ref)
		rsl := ref.Resolve(rt.Schema)
		if !reflect.DeepEqual(rsl, rt.Resolved) {
			t.Errorf("%d: resolved schema don't match, got %v, wants %v", i, rsl, rt.Resolved)
		}
	}
}

var hrefTests = []struct {
	HRef     string
	Schema   *Schema
	Order    []string
	Resolved map[string]*Schema
}{
	{
		HRef: "/app/{(%23%2Fdefinitions%2Fapp%2Fdefinitions%2Fuuid)}",
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"app": {
					Definitions: map[string]*Schema{
						"uuid": {
							Title: "Identifier",
						},
					},
				},
			},
		},
		Order: []string{"appUUID"},
		Resolved: map[string]*Schema{
			"appUUID": {
				Title: "Identifier",
			},
		},
	},
	{
		HRef: "/app/{(%23%2Fdefinitions%2Fapp%2Fdefinitions%2Fidentity)}/struct/{(%23%2Fdefinitions%2Fstruct%2Fdefinitions%2Fidentity)}",
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"app": {
					Definitions: map[string]*Schema{
						"identity": {
							Title: "App Identifier",
						},
					},
				},
				"struct": {
					Definitions: map[string]*Schema{
						"identity": {
							Title: "Struct Identifier",
						},
					},
				},
			},
		},
		Order: []string{"appIdentity", "structIdentity"},
		Resolved: map[string]*Schema{
			"appIdentity": {
				Title: "App Identifier",
			},
			"structIdentity": {
				Title: "Struct Identifier",
			},
		},
	},
}

func TestHREfResolve(t *testing.T) {
	for i, ht := range hrefTests {
		href := NewHRef(ht.HRef)
		href.Resolve(ht.Schema)
		if !reflect.DeepEqual(href.Order, ht.Order) {
			t.Errorf("%d: resolved order don't match, got %v, wants %v", i, href.Order, ht.Order)
		}
		if !reflect.DeepEqual(href.Schemas, ht.Resolved) {
			t.Errorf("%d: resolved schemas don't match, got %v, wants %v", i, href.Schemas, ht.Resolved)
		}
	}
}

var urlPatternTests = []struct {
	HRef       string
	Schema     *Schema
	URLPattern string
}{
	{
		HRef: "/app/{(%23%2Fdefinitions%2Fapp%2Fdefinitions%2Fidentifier)}",
		Schema: &Schema{
			Definitions: map[string]*Schema{
				"app": {
					Definitions: map[string]*Schema{
						"identifier": {
							Title: "Identifier",
							Type:  "integer",
						},
					},
				},
			},
		},
		URLPattern: "/app/{appIdentifier:[0-9]+}",
	},
}

func TestURLPattern(t *testing.T) {
	for i, ht := range urlPatternTests {
		href := NewHRef(ht.HRef)
		href.Resolve(ht.Schema)
		pattern := href.URLPattern()
		if pattern != ht.URLPattern {
			t.Errorf("%d: resolved pattern don't match, got %v, wants %v", i, pattern, ht.URLPattern)
		}
	}
}

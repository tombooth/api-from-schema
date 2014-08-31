package main

import (
	"strings"
	"testing"
)

func assertTemplateEquals(t *testing.T, templateFunc func(string, interface{}) (string, error), expected, name string, context interface{}) {
	output, err := templateFunc(name, context)
	if err != nil {
		t.Errorf("%v should have executed: %v", name, err)
	}

	assertEqual(t, expected, strings.TrimSpace(output))
}

func assertExecuteEquals(t *testing.T, expected, name string, context interface{}) {
	store, _ := NewTemplateStore("fixtures/templates/")
	assertTemplateEquals(t, store.Execute, expected, name, context)
}

func assertExecuteAndFormatEquals(t *testing.T, expected, name string, context interface{}) {
	store, _ := NewTemplateStore("fixtures/templates/")
	assertTemplateEquals(t, store.ExecuteAndFormat, expected, name, context)
}

func TestTemplateStoreCreation(t *testing.T) {
	_, err := NewTemplateStore("fixtures/not-a-dir/")
	if err == nil {
		t.Error("Bad directory should have blown up")
	}

	_, err = NewTemplateStore("fixtures/model-tests.json")
	if err == nil {
		t.Error("Not a directory so should have blown up")
	}

	_, err = NewTemplateStore("fixtures/templates/")
	if err != nil {
		t.Errorf("Correct directory should be fine: %v", err)
	}
}

func TestTemplateStoreExecution(t *testing.T) {
	assertExecuteEquals(t, "some-string", "echo.tmpl", "some-string")
	assertExecuteAndFormatEquals(t, "package foo", "needsFormatting.tmpl", nil)
}

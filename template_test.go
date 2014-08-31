package main

import (
	"strings"
	"testing"
)

func assertTemplateEquals(t *testing.T, templateFunc func(interface{}, ...string) (string, error), expected, context interface{}, paths ...string) {
	output, err := templateFunc(context, paths...)
	if err != nil {
		t.Errorf("%v should have executed: %v", paths, err)
	}

	assertEqual(t, expected, strings.TrimSpace(output))
}

func assertExecuteEquals(t *testing.T, expected, context interface{}, paths ...string) {
	store, _ := NewTemplateStore("fixtures/templates/")
	assertTemplateEquals(t, store.Execute, expected, context, paths...)
}

func assertExecuteAndFormatEquals(t *testing.T, expected, context interface{}, paths ...string) {
	store, _ := NewTemplateStore("fixtures/templates/")
	assertTemplateEquals(t, store.ExecuteAndFormat, expected, context, paths...)
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
	assertExecuteEquals(t, "some-string", "some-string", "echo.tmpl")
	assertExecuteAndFormatEquals(t, "package foo", nil, "needsFormatting.tmpl")
	assertExecuteEquals(t, "some-string", "some-string", "parent.tmpl", "child.tmpl")
}

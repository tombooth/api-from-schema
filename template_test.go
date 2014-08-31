package main

import (
	"strings"
	"testing"
)

func TestTemplateStoreCreation(t *testing.T) {
	_, err := NewTemplateStore("fixtures/not-a-dir/")
	if err == nil {
		t.Error("Bad directory should have blown up")
	}

	_, err = NewTemplateStore("fixtures/model-tests.json")
	if err == nil {
		t.Error("Not a directory so should have blown up")
	}

	_, err = NewTemplateStore("fixtures/templates-no-structure")
	if err == nil {
		t.Error("A directory without structure.json should errored")
	}

	_, err = NewTemplateStore("fixtures/templates/")
	if err != nil {
		t.Errorf("Correct directory should be fine: %v", err)
	}
}

func TestStructure(t *testing.T) {
	store, _ := NewTemplateStore("fixtures/templates/")

	if formattedTemplate, ok := store.Files["something.go"]; ok {
		if contents, err := formattedTemplate.Execute(nil); err != nil {
			t.Errorf("Failed to execute formatted template: %v", err)
		} else {
			assertEqual(t, "package foo", strings.TrimSpace(contents))
		}
	} else {
		t.Error("Failed to find file something.go")
	}

	if nestedTemplate, ok := store.Files["README.md"]; ok {
		if contents, err := nestedTemplate.Execute("some-value"); err != nil {
			t.Errorf("Failed to execute nested template: %v", err)
		} else {
			assertEqual(t, "some-value", strings.TrimSpace(contents))
		}
	} else {
		t.Error("Failed to find file README.md")
	}
}

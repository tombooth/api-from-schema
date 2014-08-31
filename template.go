package main

import (
	"bytes"
	"errors"
	"go/format"
	"os"
	"path"
	"text/template"
)

type TemplateStore struct {
	directory string
}

func NewTemplateStore(path string) (TemplateStore, error) {
	if finfo, err := os.Stat(path); err != nil {
		return TemplateStore{}, err
	} else if !finfo.IsDir() {
		return TemplateStore{}, errors.New("NewTemplateStore: path is not a directory")
	}

	return TemplateStore{
		directory: path,
	}, nil
}

func (store *TemplateStore) Execute(context interface{}, paths ...string) (string, error) {
	var templateOutput bytes.Buffer

	qualifiedPaths := []string{}
	for _, unqualPath := range paths {
		qualifiedPaths = append(
			qualifiedPaths,
			path.Join(store.directory, unqualPath))
	}

	tmpl, err := template.ParseFiles(qualifiedPaths...)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&templateOutput, context)
	if err != nil {
		return "", err
	}

	return templateOutput.String(), nil
}

func (store *TemplateStore) ExecuteAndFormat(context interface{}, paths ...string) (string, error) {
	stringOutput := ""

	unformattedSource, err := store.Execute(context, paths...)
	if err != nil {
		return stringOutput, err
	}

	if formattedSource, err := format.Source([]byte(unformattedSource)); err != nil {
		return stringOutput, err
	} else {
		stringOutput = string(formattedSource)
	}

	return stringOutput, nil
}

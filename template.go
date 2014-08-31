package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

type MetaTemplate struct {
	tmpl            *template.Template
	needsFormatting bool
}

type TemplateStore struct {
	Files map[string]MetaTemplate
}

func NewTemplateStore(storePath string) (TemplateStore, error) {
	if finfo, err := os.Stat(storePath); err != nil {
		return TemplateStore{}, err
	} else if !finfo.IsDir() {
		return TemplateStore{}, errors.New(fmt.Sprintf("%v is not a directory", storePath))
	}

	structureBytes, err := ioutil.ReadFile(path.Join(storePath, "structure.json"))
	if err != nil {
		return TemplateStore{}, err
	}

	files, err := parseFiles(storePath, structureBytes)
	if err != nil {
		return TemplateStore{}, err
	}

	return TemplateStore{
		Files: files,
	}, nil
}

func (meta *MetaTemplate) Execute(context interface{}) (string, error) {
	var templateOutput bytes.Buffer

	if err := meta.tmpl.Execute(&templateOutput, context); err != nil {
		return "", err
	}

	if meta.needsFormatting {
		if formattedSource, err := format.Source(templateOutput.Bytes()); err != nil {
			return "", err
		} else {
			return string(formattedSource), nil
		}
	} else {
		return templateOutput.String(), nil
	}
}

func parseFiles(basePath string, structureBytes []byte) (map[string]MetaTemplate, error) {
	var structureDescriptor map[string][]string

	err := json.Unmarshal(structureBytes, &structureDescriptor)
	if err != nil {
		return nil, err
	}

	files := make(map[string]MetaTemplate)

	for filePath, templatePaths := range structureDescriptor {
		qualifiedPaths := qualifyPaths(basePath, templatePaths)
		tmpl, err := template.ParseFiles(qualifiedPaths...)
		if err != nil {
			return nil, err
		}
		files[filePath] = MetaTemplate{
			tmpl:            tmpl,
			needsFormatting: strings.HasSuffix(filePath, ".go"),
		}
	}

	return files, nil
}

func qualifyPaths(basePath string, unqualifiedPaths []string) []string {
	qualifiedPaths := []string{}
	for _, unqualifiedPath := range unqualifiedPaths {
		qualifiedPaths = append(
			qualifiedPaths,
			path.Join(basePath, unqualifiedPath))
	}
	return qualifiedPaths
}

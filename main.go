package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/docopt/docopt-go"
	"github.com/tombooth/api-from-schema/schematic"
)

type Project struct {
	Name     string
	Username string
}

func ProjectFromArguments(arguments map[string]interface{}) (Project, error) {
	username, ok := arguments["--username"].(string)
	if !ok {
		if currentUser, err := user.Current(); err != nil {
			return Project{}, err
		} else {
			username = currentUser.Username
		}
	}

	name, ok := arguments["--name"].(string)
	if !ok {
		name = path.Base(arguments["<project_path>"].(string))
	}

	return Project{
		Name:     name,
		Username: username,
	}, nil
}

func main() {
	usage := `Api from schema

Usage:
  api-from-schema [options] <json_schema> <project_path>
  api-from-schema -h | --help
  api-from-schema --version

Options:
  -h --help               Show this screen.
  --version               Show version.
  --templates=<path>      Path to template directory [default: templates/]
  --username=<username>   Github username, defaults to system username
  --name=<name>           Project name, defaults to directory name in project path
`

	arguments, _ := docopt.Parse(usage, nil, true, "Api from schema 0.1.0", false)
	project, err := ProjectFromArguments(arguments)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load project configuration: %v", err)
		return
	}

	schemaPath := arguments["<json_schema>"].(string)
	templateStore, err := NewTemplateStore(arguments["--templates"].(string))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load templates: %v", err)
		return
	}

	projectPath, err := ensureDirectory(arguments["<project_path>"].(string))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ensure project directory: %v", err)
		return
	}

	if apiSchema, err := schema.ParseSchema(schemaPath); err == nil {
		apiSchema.Resolve(nil)
		models := ModelsFromSchema(apiSchema)

		context := struct {
			Project Project
			Models  []Model
		}{
			Project: project,
			Models:  models,
		}

		for filePath, fileTemplate := range templateStore.Files {
			fullPath := path.Join(projectPath, filePath)
			if contents, err := fileTemplate.Execute(context); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to execute template: %v", err)
			} else {
				writeFile(fullPath, contents)
			}
		}
	}
}

func question(format string, a ...interface{}) error {
	var answer string

	fmt.Printf(
		"%v, are you sure you want to continue? (y/n) ",
		fmt.Sprintf(format, a...),
	)
	if _, err := fmt.Scanf("%s", &answer); err != nil {
		return err
	} else if answer != "y" {
		return errors.New(fmt.Sprintf("Aborted with answer %v", answer))
	}

	return nil
}

func ensureDirectory(projectPath string) (string, error) {
	cleanedPath := path.Clean(projectPath)

	if finfo, err := os.Stat(cleanedPath); err != nil && !os.IsNotExist(err) {
		return "", err
	} else if finfo != nil && !finfo.IsDir() {
		return "", errors.New(fmt.Sprintf("%v is a file not a directory", projectPath))
	} else if finfo != nil && finfo.IsDir() {
		if err := question("%v exists", projectPath); err != nil {
			return "", err
		}
	} else {
		if err := question("Creating directory %v", projectPath); err != nil {
			return "", err
		}
		if err := os.MkdirAll(cleanedPath, 0755); err != nil {
			return "", err
		}
	}

	return cleanedPath, nil
}

func writeFile(filePath, contents string) {
	fmt.Printf("Writing out to %v\n", filePath)
	containingDirectory := path.Dir(filePath)
	if err := os.MkdirAll(containingDirectory, 0755); err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to create containg directory %v: %v\n", containingDirectory, err)
	} else {
		if file, err := os.Create(filePath); err != nil {
			fmt.Fprintf(os.Stderr,
				"Failed to open %v to write into: %v\n", filePath, err)
		} else {
			fmt.Fprint(file, contents)
		}
	}
}

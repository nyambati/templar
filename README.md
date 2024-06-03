# templar
Golang tool for processing scaffolding templates. Its inspired by cookiecutter

## Usage

### Cli
Templar can be used as a cli.
```
$ go install github.com/nyambati/templar
$ templar generate -c templar.yaml
```
Configuration
```yaml
template_dir: testdata/templates
output_dir: testdata/out
overwrite: true
variables:
  vertical: test
  environment: staging
  region: us-east-1
  account_id: 123456789
  account_name: testaccount

```

### Package

You can also use this tool as a package within your Go applications.

```go
package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/nyambati/templar/pkg/generator"
)

type Vars map[string]string

func (v *Vars) Validate() error {
	// validate vars
	return nil
}

func main() {
	templateDir := "./templates"
	outputDir := "./output"

	vars := Vars{
		"vertical":     "test",
		"environment":  "dev",
		"region":       "us-east-1",
		"account_id":   "123456789",
		"account_name": "testaccount",
	}

	// Prompt user for input
	answers := make(map[string]interface{})
	for k, v := range vars {
		answers[k] = v
	}

	err := prompt(vars, &answers)
	if err != nil {
		panic(err)
	}

	generator, err := generator.New(
		templateDir,
		outputDir,
		true,
		&vars,
	)
	if err != nil {
		panic(err)
	}

	if err := generator.Generate(); err != nil {
		panic(err)
	}
}

func prompt(vars map[string]string, answers *map[string]interface{}) error {
	questions := make([]*survey.Question, 0, len(vars))
	for key, value := range vars {
		questions = append(questions, &survey.Question{
			Name:   key,
			Prompt: &survey.Input{Message: key, Default: value},
		})
	}
	return survey.Ask(questions, answers)
}

```
## Templating

This tool leverages Go's text/template package for templating, enhanced with the Sprig library for additional template functions.

### Using Go Templates
Templates are defined using Go's text/template syntax.
```txr
templates/{{project_name}}/file.hcl
```
To enhance the templating exprecience, we have used sprig to extend templating functions. Sprig provides over 70 additional template functions that you can use in your templates.

```
# lowercase
templates/mytemplate/folder1/{{ProjectName | lower}}/file1.hcl

# upper
hello = "{{.ProjectName | upper}}"

```


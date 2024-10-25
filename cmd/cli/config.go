package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Vars map[string]string

type Config struct {
	TemplatesDir string
	OutputDir    string
	Overwrite    bool
	Vars         Vars
}

func (vars *Vars) Validate() error {
	required := []string{
		"vertical",
		"environment",
		"region",
		"account_id",
		"account_name",
	}
	errs := []error{}
	for _, value := range required {
		if _, ok := (*vars)[value]; !ok {
			errs = append(errs, fmt.Errorf("%s is required", value))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func NewConfig(path string) Config {
	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("%s/templar.yaml", path))
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return Config{
		TemplatesDir: v.GetString("template_dir"),
		OutputDir:    v.GetString("output_dir"),
		Overwrite:    v.GetBool("overwrite"),
		Vars:         v.GetStringMapString("variables"),
	}
}

func Prompt(vars map[string]string, answers *map[string]interface{}) error {
	questions := make([]*survey.Question, 0, len(vars))
	keys := maps.Keys(vars)
	slices.Sort(keys)
	for _, key := range keys {
		questions = append(questions, &survey.Question{
			Name: key,
			Prompt: &survey.Input{
				Message: key,
				Default: vars[key],
			},
		})
	}
	return survey.Ask(questions, answers)
}

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

func (c *Vars) Validate() error {
	required := []string{
		"vertical",
		"environment",
		"region",
		"account_id",
		"account_name",
	}
	errs := []error{}
	for _, v := range required {
		if _, ok := (*c)[v]; !ok {
			errs = append(errs, fmt.Errorf("%s is required", v))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func NewConfig(path string) Config {
	if path == "" {
		path = "."
	}
	v := viper.New()
	v.SetConfigFile("templar.yaml")
	v.AddConfigPath(path)
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

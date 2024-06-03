package main

import (
	"github.com/nyambati/templar/pkg/generator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var answers = make(map[string]interface{})
var config Config

var rootCmd = &cobra.Command{
	Use: "templar",
}

var generatorCmd = &cobra.Command{
	Use: "generate",
	PreRun: func(cmd *cobra.Command, args []string) {
		config = NewConfig(".")
		for k, v := range config.Vars {
			answers[k] = v
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := Prompt(config.Vars, &answers); err != nil {
			logrus.Fatal(err)
		}

		// convert answers to Vars
		for k, v := range answers {
			config.Vars[k] = v.(string)
		}

		generator, err := generator.New(
			config.TemplatesDir,
			config.OutputDir,
			config.Overwrite,
			&config.Vars,
		)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := generator.Generate(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(generatorCmd)
	rootCmd.Execute()
}

package main

import (
	"github.com/nyambati/templar/pkg/generator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var answers = make(map[string]interface{})
var config Config
var configPath string

var rootCmd = &cobra.Command{
	Use: "templar",
}

var generatorCmd = &cobra.Command{
	Use: "generate",
	PreRun: func(cmd *cobra.Command, args []string) {
		config = NewConfig(configPath)
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

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.TemplatesDir, "template-dir", "t", "", "Templates directory")
	rootCmd.PersistentFlags().StringVarP(&config.OutputDir, "output-dir", "o", "", "Output directory")
	rootCmd.PersistentFlags().BoolVarP(&config.Overwrite, "overwrite", "w", false, "Overwrite existing files")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".", "Config file path")
}

package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Generator struct {
	InputDir  string
	OutputDir string
	Overwrite bool
	Vars      Vars
}

// GenerateProject creates a project from a template directory.
func (g *Generator) Generate() error {
	err := filepath.Walk(g.InputDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Determine the new path in the output directory
		newPath, err := g.parsePath(srcPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return os.MkdirAll(newPath, os.ModePerm)
		}

		return g.parseFile(srcPath, newPath)
	})
	return err
}

func (g *Generator) parseFile(file, path string) error {
	if !strings.Contains(file, ".hcl") {
		return nil
	}
	buffer, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tmpl, err := template.New(file).Funcs(sprig.FuncMap()).Parse(string(buffer))
	if err != nil {
		return err
	}

	// if file already exists and overwrite is not set, skip
	if _, err := os.Stat(path); err == nil && !g.Overwrite {
		return nil
	}

	outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return tmpl.Execute(outFile, g.Vars)
}

func (g *Generator) parsePath(srcPath string) (string, error) {
	relPath, err := filepath.Rel(g.InputDir, srcPath)
	if err != nil {
		return "", err
	}

	// Process template for directory names
	tmpl, err := template.New("path").Funcs(sprig.FuncMap()).Parse(relPath)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := tmpl.Execute(&builder, g.Vars); err != nil {
		return "", err
	}
	return filepath.Join(g.OutputDir, builder.String()), nil
}

func New(inputDir, outputDir string, overwrite bool, vars Vars) (GeneratorInterface, error) {
	// validate vars
	if err := vars.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate vars: %w", err)
	}
	return &Generator{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Overwrite: overwrite,
		Vars:      vars,
	}, nil
}

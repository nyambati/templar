package generator_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/nyambati/templar/pkg/generator"
	"github.com/stretchr/testify/assert"
)

type Vars map[string]string

func (vars *Vars) Validate() error {
	required := []string{
		"vertical",
		"environment",
		"region",
		"account_id",
		"account_name",
	}
	errs := []error{}
	for _, v := range required {
		if _, ok := (*vars)[v]; !ok {
			errs = append(errs, fmt.Errorf("%s is required", v))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func TestGenerator_Generate(t *testing.T) {
	type fields struct {
		InputDir  string
		OutputDir string
		Overwrite bool
		Vars      generator.Vars
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "WithValidVars",
			fields: fields{
				InputDir:  "../../testdata/templates",
				OutputDir: "../../testdata/out",
				Overwrite: false,
				Vars: &Vars{
					"vertical":     "dp-cloud-infra",
					"environment":  "dev",
					"region":       "us-east-1",
					"account_id":   "123456789",
					"account_name": "dp-cloud-infra",
				},
			},
			wantErr: false,
		},
		{
			name: "WithInvalidVars",
			fields: fields{
				InputDir:  "../../testdata/templates",
				OutputDir: "../../testdata/out",
				Overwrite: false,
				Vars: &Vars{
					"vertical":    "dp-cloud-infra",
					"environment": "dev",
					"region":      "us-east-1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		path, _ := filepath.Abs(tt.fields.InputDir)
		t.Run(tt.name, func(t *testing.T) {
			g, err := generator.New(path, tt.fields.OutputDir, tt.fields.Overwrite, tt.fields.Vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.New() error = %v", err)
			} else if g == nil {
				return
			}
			if err := g.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			accountDir := "../../testdata/out/test/production/"
			assert.DirExists(t, accountDir)
			assert.DirExists(t, accountDir+"/test/us-east-1")
			assert.DirExists(t, accountDir+"/test/_global")
			assert.FileExists(t, accountDir+"/test/_global/region.hcl")
			assert.FileExists(t, accountDir+"/test/us-east-1/region.hcl")
			assert.FileExists(t, accountDir+"/environment.hcl")
		})
	}
}

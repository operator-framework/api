package internal

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/operator-framework/api/pkg/validation/errors"

	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func TestValidateCSV(t *testing.T) {
	cases := []struct {
		validatorFuncTest
		csvPath string
	}{
		{
			validatorFuncTest{
				description: "successfully validated",
			},
			filepath.Join("testdata", "correct.csv.yaml"),
		},
		{
			validatorFuncTest{
				description: "data type mismatch",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidParse(
						`converting bundle CSV "etcdoperator.v0.9.0"`,
						"json: cannot unmarshal string into Go struct field ClusterServiceVersionSpec.maintainers of type []v1alpha1.Maintainer"),
				},
			},
			filepath.Join("testdata", "dataTypeMismatch.csv.yaml"),
		},
	}
	for _, c := range cases {
		b, err := ioutil.ReadFile(c.csvPath)
		if err != nil {
			t.Fatalf("Error reading CSV path %s: %v", c.csvPath, err)
		}
		csv := registry.ClusterServiceVersion{}
		if err = yaml.Unmarshal(b, &csv); err != nil {
			t.Fatalf("Error unmarshalling CSV at path %s: %v", c.csvPath, err)
		}
		result := validateCSVRegistry(&csv)
		c.check(t, result)
	}
}

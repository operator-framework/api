package internal

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/api/pkg/validation/errors"

	"github.com/ghodss/yaml"
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
				description: "invalid install modes",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidCSV("install modes not found", "etcdoperator.v0.9.0"),
				},
			},
			filepath.Join("testdata", "noInstallMode.csv.yaml"),
		},
	}
	for _, c := range cases {
		b, err := ioutil.ReadFile(c.csvPath)
		if err != nil {
			t.Fatalf("Error reading CSV path %s: %v", c.csvPath, err)
		}
		csv := operatorsv1alpha1.ClusterServiceVersion{}
		if err = yaml.Unmarshal(b, &csv); err != nil {
			t.Fatalf("Error unmarshalling CSV at path %s: %v", c.csvPath, err)
		}
		result := validateCSV(&csv)
		c.check(t, result)
	}
}

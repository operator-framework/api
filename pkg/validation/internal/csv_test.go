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
		{
			validatorFuncTest{
				description: "valid install modes when dealing with conversionCRDs",
				wantErr:     false,
			},
			filepath.Join("testdata", "correct.csv.with.conversion.webhook.yaml"),
		},
		{
			validatorFuncTest{
				description: "invalid install modes when dealing with conversionCRDs",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidCSV("only AllNamespaces InstallModeType is supported when conversionCRDs is present", "etcdoperator.v0.9.0"),
				},
			},
			filepath.Join("testdata", "incorrect.csv.with.conversion.webhook.yaml"),
		},
		{
			validatorFuncTest{
				description: "invalid annotation name for csv",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrFailedValidation("provided annotation olm.skiprange uses wrong case and should be olm.skipRange instead", "etcdoperator.v0.9.0"),
					errors.ErrFailedValidation("provided annotation olm.operatorgroup uses wrong case and should be olm.operatorGroup instead", "etcdoperator.v0.9.0"),
					errors.ErrFailedValidation("provided annotation olm.operatornamespace uses wrong case and should be olm.operatorNamespace instead", "etcdoperator.v0.9.0"),
				},
			},
			filepath.Join("testdata", "badAnnotationNames.csv.yaml"),
		},
		{
			validatorFuncTest{
				description: "csv with name over 63 characters limit",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidCSV(`metadata.name "someoperatorwithanextremelylongnamethatmakenosensewhatsoever.v999.999.999" is invalid: must be no more than 63 characters`, "someoperatorwithanextremelylongnamethatmakenosensewhatsoever.v999.999.999"),
				},
			},
			filepath.Join("testdata", "badName.csv.yaml"),
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

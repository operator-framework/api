package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/operator-framework/api/pkg/validation/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
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
		{
			validatorFuncTest{
				description: "should fail when alm-examples is pretty format and is invalid",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidParse("invalid example", "invalid character at 176\n [{\"apiVersion\":\"local.storage.openshift.io/v1\",\"kind\":\"LocalVolume\",\"metadata\":{\"name\":\"example\"},\"spec\":{\"storageClassDevices\":[{\"devicePaths\":[\"/dev/disk/by-id/ata-crucial\",]<--(see the invalid character)"),
				},
			},
			filepath.Join("testdata", "invalid.alm-examples.csv.yaml"),
		},
		{
			validatorFuncTest{
				description: "should not fail when alm-examples is not informed",
				wantWarn:    true,
				errors: []errors.Error{
					errors.WarnInvalidOperation("provided API should have an example annotation", schema.GroupVersionKind{Group: "etcd.database.coreos.com", Version: "v1beta2", Kind: "EtcdCluster"}),
				},
			},
			filepath.Join("testdata", "correct.csv.empty.example.yaml"),
		},
		{
			validatorFuncTest{
				description: "should warn when olm.properties are defined in the annotations",
				wantWarn:    true,
				errors: []errors.Error{
					errors.WarnPropertiesAnnotationUsed(
						fmt.Sprintf(
							"found %s annotation, please define these properties in metadata/properties.yaml instead",
							olmpropertiesAnnotation,
						),
					),
				},
			},
			filepath.Join("testdata", "correct.csv.olm.properties.annotation.yaml"),
		},
		{
			validatorFuncTest{
				description: "should fail when spec.minKubeVersion is not in semantic version format",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidCSV(`csv.Spec.MinKubeVersion has an invalid value: 1.21`, "test-operator.v0.0.1"),
				},
			},
			filepath.Join("testdata", "invalid_min_kube_version.csv.yaml"),
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

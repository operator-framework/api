package internal

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/stretchr/testify/require"
)

func TestValidateBundle(t *testing.T) {
	var table = []struct {
		description string
		directory   string
		hasError    bool
		errString   string
	}{
		{
			description: "registryv1 valid bundle",
			directory:   "./testdata/valid_bundle",
			hasError:    false,
		},
		{
			description: "registryv1 valid bundle with multiple versions in CRD",
			directory:   "./testdata/valid_bundle_2",
			hasError:    false,
		},
		{
			description: "registryv1 invalid bundle without CRD etcdclusters v1beta2 in bundle",
			directory:   "./testdata/invalid_bundle",
			hasError:    true,
			errString:   `owned CRD "etcd.database.coreos.com/v1beta2, Kind=EtcdCluster" not found in bundle`,
		},
		{
			description: "registryv1 invalid bundle without CRD etcdclusters v1beta2 in CSV",
			directory:   "./testdata/invalid_bundle_2",
			hasError:    true,
			errString:   `CRD "etcd.database.coreos.com/v1beta2, Kind=EtcdCluster" is present in bundle "etcdoperator.v0.9.4" but not defined in CSV`,
		},
		{
			description: "registryv1 invalid bundle with duplicate CRD etcdclusters in bundle",
			directory:   "./testdata/invalid_bundle_3",
			hasError:    true,
			errString:   `duplicate CRD "test.example.com/v1alpha1, Kind=Test" in bundle "test-operator.v0.0.1"`,
		},
		{
			description: "invalid bundle service account can't match sa in csv",
			directory:   "./testdata/invalid_bundle_sa",
			hasError:    true,
			errString:   `invalid service account found in bundle. sa name cannot match service account defined for deployment spec in CSV`,
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.directory)
			require.NoError(t, err)

			results := BundleValidator.Validate(bundle)

			require.Greater(t, len(results), 0)
			if tt.hasError {
				require.True(t, results[0].HasError(), "found no error when an error was expected")
				require.Contains(t, results[0].Errors[0].Error(), tt.errString)
			} else {
				require.False(t, results[0].HasError(), "found error when an error was not expected")
			}
		})
	}
}

func TestValidateServiceAccount(t *testing.T) {
	csvWithSAs := func(saNames ...string) *v1alpha1.ClusterServiceVersion {
		csv := &v1alpha1.ClusterServiceVersion{}
		depSpecs := make([]v1alpha1.StrategyDeploymentSpec, len(saNames))
		for i, saName := range saNames {
			depSpecs[i].Spec.Template.Spec.ServiceAccountName = saName
		}
		csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs = depSpecs
		return csv
	}

	var table = []struct {
		description string
		bundle      *manifests.Bundle
		hasError    bool
		errString   string
	}{
		{
			description: "an object with the same name as the service account",
			bundle: &manifests.Bundle{
				CSV: csvWithSAs("foo"),
				Objects: []*unstructured.Unstructured{
					{Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name": "foo",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"serviceAccountName": "foo",
								},
							},
						},
					}},
				},
			},
			hasError: false,
		},
		{
			description: "service account included in both CSV and bundle",
			bundle: &manifests.Bundle{
				CSV: csvWithSAs("foo"),
				Objects: []*unstructured.Unstructured{
					{Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name": "foo",
						},
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"serviceAccountName": "foo",
								},
							},
						},
					}},
					{Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ServiceAccount",
						"metadata": map[string]interface{}{
							"name": "foo",
						},
					}},
				},
			},
			hasError:  true,
			errString: `invalid service account found in bundle. sa name cannot match service account defined for deployment spec in CSV`,
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Validate the bundle object
			results := BundleValidator.Validate(tt.bundle)

			require.Greater(t, len(results), 0)
			if tt.hasError {
				require.True(t, results[0].HasError(), "found no error when an error was expected")
				require.Contains(t, results[0].Errors[0].Error(), tt.errString)
			} else {
				require.False(t, results[0].HasError(), "found error when an error was not expected")
			}
		})
	}
}

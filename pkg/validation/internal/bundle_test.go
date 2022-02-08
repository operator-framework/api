package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/api/pkg/validation/errors"
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
			errString:   `invalid service account found in bundle. This service account etcd-operator in your bundle is not valid, because a service account with the same name was already specified in your CSV. If this was unintentional, please remove the service account manifest from your bundle. If it was intentional to specify a separate service account, please rename the SA in either the bundle manifest or the CSV.`,
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
			errString: `invalid service account found in bundle. This service account foo in your bundle is not valid, because a service account with the same name was already specified in your CSV. If this was unintentional, please remove the service account manifest from your bundle. If it was intentional to specify a separate service account, please rename the SA in either the bundle manifest or the CSV.`,
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

func TestBundleSize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name        string
		args        args
		wantError   bool
		wantWarning bool
		errStrings  []string
		warnStrings []string
	}{
		{
			name: "should pass when the size is not bigger or closer of the limit",
			args: args{
				size: int64(max_bundle_size / 2),
			},
		},
		{
			name: "should warn when the size is closer of the limit",
			args: args{
				size: int64(max_bundle_size - 10),
			},
			wantWarning: true,
			warnStrings: []string{"Warning: : nearing maximum bundle compressed size with gzip: size=~3 MegaByte, max=4 MegaByte"},
		},
		{
			name:        "should warn when is not possible to check the size",
			wantWarning: true,
			warnStrings: []string{"Warning: : unable to check the bundle size"},
		},
		{
			name: "should raise an error when the size is bigger than the limit",
			args: args{
				size: int64(2 * max_bundle_size),
			},
			wantError:  true,
			errStrings: []string{"Error: : maximum bundle compressed size with gzip size exceeded: size=~8 MegaByte, max=4 MegaByte"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle := &manifests.Bundle{
				CompressedSize: &tt.args.size,
			}
			result := validateBundleSize(bundle)

			var warns, errs []errors.Error
			for _, r := range result {
				if r.Level == errors.LevelWarn {
					warns = append(warns, r)
				} else if r.Level == errors.LevelError {
					errs = append(errs, r)
				}
			}
			require.Equal(t, tt.wantWarning, len(warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(warns))
				for _, w := range warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(errs) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(errs))
				for _, err := range errs {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

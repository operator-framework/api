package internal

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"

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
	}

	for _, tt := range table {
		// Validate the bundle object
		bundle, err := manifests.GetBundleFromDir(tt.directory)
		require.NoError(t, err)

		results := BundleValidator.Validate(bundle)

		if len(results) > 0 {
			require.Equal(t, results[0].HasError(), tt.hasError)
			if results[0].HasError() {
				require.Contains(t, results[0].Errors[0].Error(), tt.errString)
			}
		}
	}
}

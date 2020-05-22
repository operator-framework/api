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
			description: "registryv1 bundle/valid bundle",
			directory:   "./testdata/valid_bundle",
			hasError:    false,
		},
		{
			description: "registryv1 bundle/valid bundle",
			directory:   "./testdata/invalid_bundle",
			hasError:    true,
			errString:   `owned CRD "etcdclusters.etcd.database.coreos.com" not found in bundle`,
		},
		{
			description: "registryv1 bundle/valid bundle",
			directory:   "./testdata/invalid_bundle_2",
			hasError:    true,
			errString:   `owned CRD "etcdclusters.etcd.database.coreos.com" is present in bundle "etcdoperator.v0.9.4" but not defined in CSV`,
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

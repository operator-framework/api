package internal

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"

	"github.com/stretchr/testify/require"
)

func TestValidateCapabilities(t *testing.T) {
	var table = []struct {
		description string
		directory   string
		hasError    bool
		errStrings  []string
	}{
		{
			description: "registryv1 bundle/valid bundle",
			directory:   "./testdata/valid_bundle",
			hasError:    false,
		},
		{
			description: "registryv1 bundle/invald bundle operatorhubio",
			directory:   "./testdata/invalid_bundle_operatorhub",
			hasError:    true,
			errStrings: []string{
				`Error: Value : (etcdoperator.v0.9.4) csv.Metadata.Annotations.Capabilities "Installs and stuff" is not a valid capabilities level`,
			},
		},
	}

	for _, tt := range table {
		// Validate the bundle object
		bundle, err := manifests.GetBundleFromDir(tt.directory)
		require.NoError(t, err)

		results := StandardCapabilitiesValidator.Validate(bundle)

		if len(results) > 0 {
			require.Equal(t, results[0].HasError(), tt.hasError)
			if results[0].HasError() {
				require.Equal(t, len(tt.errStrings), len(results[0].Errors))

				for _, err := range results[0].Errors {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		}
	}
}

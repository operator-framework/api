package internal

import (
	"io/ioutil"
	"testing"

	"github.com/operator-framework/api/pkg/validation/errors"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/ghodss/yaml"
)

func TestValidateCRD(t *testing.T) {
	var table = []struct {
		description string
		directory   string
		version     string
		hasError    bool
		errString   string
	}{
		{
			description: "registryv1 bundle/valid crd",
			directory:   "./testdata/v1beta1.crd.yaml",
			version:     "v1beta1",
			hasError:    false,
			errString:   "",
		},
		{
			description: "registryv1 bundle/invalid crd",
			directory:   "./testdata/duplicateVersions.crd.yaml",
			version:     "v1beta1",
			hasError:    true,
			errString:   "must contain unique version names",
		},
		{
			description: "registryv1 bundle/invalid crd",
			directory:   "./testdata/v1.crd.yaml",
			version:     "v1",
			hasError:    false,
			errString:   "",
		},
		{
			description: "registryv1 bundle/invalid crd",
			directory:   "./testdata/deprecatedVersion.crd.yaml",
			version:     "v1",
			hasError:    true,
			errString:   "must have exactly one version marked as storage version",
		},
	}

	for _, tt := range table {
		b, err := ioutil.ReadFile(tt.directory)
		if err != nil {
			t.Fatalf("Error reading CRD path %s: %v", tt.directory, err)
		}

		results := []errors.ManifestResult{}
		switch tt.version {
		case "v1":
			crd := &v1.CustomResourceDefinition{}
			if err = yaml.Unmarshal(b, crd); err != nil {
				t.Fatalf("Error unmarshalling CRD at path %s: %v", tt.directory, err)
			}
			results = CRDValidator.Validate(crd)
		default:
			crd := &v1beta1.CustomResourceDefinition{}
			if err = yaml.Unmarshal(b, crd); err != nil {
				t.Fatalf("Error unmarshalling CRD at path %s: %v", tt.directory, err)
			}
			results = CRDValidator.Validate(crd)
		}

		if len(results) > 0 {
			require.Equal(t, results[0].HasError(), tt.hasError)
			if results[0].HasError() {
				require.Contains(t, results[0].Errors[0].Error(), tt.errString)
			}
		}
	}
}

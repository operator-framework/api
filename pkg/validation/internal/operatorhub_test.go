package internal

import (
	"os"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"

	"github.com/stretchr/testify/require"
)

func TestValidateBundleOperatorHub(t *testing.T) {
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
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Provider.Name not specified`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Maintainers elements should contain both name and email`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Maintainers email invalidemail is invalid: mail: missing '@' or angle-addr`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Links elements should contain both name and url`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Links url https//coreos.com/operators/etcd/docs/latest/ is invalid: parse "https//coreos.com/operators/etcd/docs/latest/": invalid URI for request`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Metadata.Annotations.Capabilities Installs and stuff is not a valid capabilities level`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Icon should only have one element`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Metadata.Annotations.Categories Magic is not a valid category`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Version is not set`,
			},
		},
	}

	for _, tt := range table {
		// Validate the bundle object
		bundle, err := manifests.GetBundleFromDir(tt.directory)
		require.NoError(t, err)

		results := OperatorHubValidator.Validate(bundle)

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

func TestCustomCategories(t *testing.T) {
	var table = []struct {
		description string
		directory   string
		hasError    bool
		errStrings  []string
		custom      bool
	}{
		{
			description: "valid bundle custom categories",
			directory:   "./testdata/valid_bundle_custom_categories",
			hasError:    false,
			custom:      true,
		},
		{
			description: "valid bundle standard categories",
			directory:   "./testdata/valid_bundle",
			hasError:    false,
			custom:      false,
		},
	}

	for _, tt := range table {
		if tt.custom {
			os.Setenv("OPERATOR_BUNDLE_CATEGORIES", "./testdata/categories.json")
		} else {
			os.Setenv("OPERATOR_BUNDLE_CATEGORIES", "")
		}

		// Validate the bundle object
		bundle, err := manifests.GetBundleFromDir(tt.directory)
		require.NoError(t, err)

		results := OperatorHubValidator.Validate(bundle)

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

func TestExtractCategories(t *testing.T) {
	path := "./testdata/categories.json"
	categories, err := extractCategories(path)
	if err != nil {
		t.Fatalf("extracting categories.json: %s", err)
	}

	expected := map[string]struct{}{
		"Cloud Pak":      {},
		"Registry":       {},
		"MyCoolThing":    {},
		"This/Or & That": {},
	}

	for key := range categories {
		if _, ok := expected[key]; !ok {
			t.Fatalf("did not find key %s", key)
		}
	}
}

func TestWarnNoIcon(t *testing.T) {
	var table = []struct {
		description string
		directory   string
		hasError    bool
		errStrings  []string
		warnStrings []string
	}{
		{
			description: "valid bundle no icon",
			directory:   "./testdata/valid_bundle_no_icon",
			hasError:    false,
			errStrings:  []string{},
			warnStrings: []string{"Warning: Value : (etcdoperator.v0.9.4) csv.Spec.Icon not specified"},
		},
	}

	for _, tt := range table {
		// Validate the bundle object
		bundle, err := manifests.GetBundleFromDir(tt.directory)
		require.NoError(t, err)

		results := OperatorHubValidator.Validate(bundle)

		if len(results) > 0 {
			require.Equal(t, tt.hasError, results[0].HasError())
			if results[0].HasError() {
				require.Equal(t, len(tt.errStrings), len(results[0].Errors))
				for _, err := range results[0].Errors {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
			if results[0].HasWarn() {
				require.Equal(t, len(tt.warnStrings), len(results[0].Warnings))
				for _, warn := range results[0].Warnings {
					warnString := warn.Error()
					require.Contains(t, tt.warnStrings, warnString)
				}
			}
		}
	}
}

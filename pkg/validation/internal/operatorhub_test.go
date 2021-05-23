package internal

import (
	"os"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

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
				`Error: Value : (etcdoperator.v0.9.4) csv.Metadata.Annotations["categories"] value Magic is not in the set of default categories`,
				`Error: Value : (etcdoperator.v0.9.4) csv.Spec.Version must be set`,
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

func TestCheckSpecIcon(t *testing.T) {
	validIcon := v1alpha1.Icon{
		Data:      "iVBORw0KGgoAAAANSUhEUgAAAOEAAADZCAYAAADWmle6AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAGXRFWHRTb2Z0d2Fy",
		MediaType: "image/png",
	}
	invalidIcon := v1alpha1.Icon{
		MediaType: "image/png",
	}

	invalidMediaTypeIcon := validIcon
	invalidMediaTypeIcon.MediaType = "invalid"

	type args struct {
		icon []v1alpha1.Icon
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
			name: "should work with a valid value",
			args: args{icon: []v1alpha1.Icon{validIcon}},
		},
		{
			name:        "should return an warning when the icon is not provided",
			wantWarning: true,
			warnStrings: []string{"csv.Spec.Icon not specified"},
		},
		{
			name:       "should fail when the data informed for the icon is invalid",
			args:       args{icon: []v1alpha1.Icon{invalidIcon}},
			wantError:  true,
			errStrings: []string{"csv.Spec.Icon elements should contain both data and mediatype"},
		},
		{
			name:       "should fail when the data informed has not a valid MediaType",
			args:       args{icon: []v1alpha1.Icon{invalidMediaTypeIcon}},
			wantError:  true,
			errStrings: []string{"csv.Spec.Icon invalid does not have a valid mediatype"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := v1alpha1.ClusterServiceVersion{}
			csv.Spec.Icon = tt.args.icon
			checks := CSVChecks{csv: csv, errs: []error{}, warns: []error{}}

			result := checkSpecIcon(checks)

			require.Equal(t, tt.wantWarning, len(result.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(result.warns))
				for _, w := range result.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(result.errs) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(result.errs))
				for _, err := range result.errs {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

func TestCheckBundleName(t *testing.T) {
	type args struct {
		bundleName string
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
			name: "should work with valid bundle name",
			args: args{bundleName: "memcached-operator.v0.9.2"},
		},
		{
			name:        "should return a warning when the bundle name is not following the convention",
			args:        args{bundleName: "memcached-operator0.9.2"},
			wantWarning: true,
			warnStrings: []string{"csv.metadata.Name memcached-operator0.9.2 is not following the recommended naming " +
				"convention: <operator-name>.v<semver> e.g. memcached-operator.v0.0.1"},
		},
		{
			name:        "should return a warning when the bundle name version is not following semver",
			args:        args{bundleName: "memcached-operator.v1"},
			wantWarning: true,
			warnStrings: []string{"csv.metadata.Name memcached-operator.v1 is not following the versioning convention " +
				"(MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/"},
		},
		{
			name:        "should return a warning when the bundle name version is not following semver",
			args:        args{bundleName: "memcached-operator.v1"},
			wantWarning: true,
			warnStrings: []string{"csv.metadata.Name memcached-operator.v1 is not following the " +
				"versioning convention (MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/"},
		},
		{
			name:        "should return a warning when the bundle name version is not following semver",
			args:        args{bundleName: "memcached-operator.v1--1.0"},
			wantWarning: true,
			warnStrings: []string{"csv.metadata.Name memcached-operator.v1--1.0 is not following the " +
				"versioning convention (MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/"},
		},
		{
			name:        "should return a warning when the bundle name version is not following semver",
			args:        args{bundleName: "memcached-operator.v1.3"},
			wantWarning: true,
			warnStrings: []string{"csv.metadata.Name memcached-operator.v1.3 is not following the " +
				"versioning convention (MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			csv := v1alpha1.ClusterServiceVersion{}
			csv.Name = tt.args.bundleName
			checks := CSVChecks{csv: csv, errs: []error{}, warns: []error{}}
			result := checkBundleName(checks)

			require.Equal(t, tt.wantWarning, len(result.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(result.warns))
				for _, w := range result.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(result.errs) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(result.errs))
				for _, err := range result.errs {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

func TestCheckSpecMinKubeVersion(t *testing.T) {
	type args struct {
		minKubeVersion string
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
			name: "should work with a valid value",
			args: args{minKubeVersion: "1.16"},
		},
		{
			name:        "should return a warning when the minKubeVersion is not informed ",
			args:        args{minKubeVersion: ""},
			wantWarning: true,
			warnStrings: []string{minKubeVersionWarnMessage},
		},
		{
			name:       "should fail when an invalid value is informed",
			args:       args{minKubeVersion: "alpha1"},
			wantError:  true,
			errStrings: []string{"csv.Spec.MinKubeVersion has an invalid value: alpha1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := v1alpha1.ClusterServiceVersion{}
			csv.Spec.MinKubeVersion = tt.args.minKubeVersion
			checks := CSVChecks{csv: csv, errs: []error{}, warns: []error{}}

			result := checkSpecMinKubeVersion(checks)

			require.Equal(t, tt.wantWarning, len(result.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(result.warns))
				for _, w := range result.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(result.errs) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(result.errs))
				for _, err := range result.errs {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

func TestValidateHubDeprecatedAPIS(t *testing.T) {
	type args struct {
		minKubeVersion string
		k8sVersion     string
		directory      string
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
			name: "should not return error or warning when the k8sVersion is <= 1.15",
			args: args{
				k8sVersion:     "1.15",
				minKubeVersion: "",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantWarning: true,
			warnStrings: []string{"checking APIs against Kubernetes version : 1.15"},
		},
		{
			name: "should return a warning when has the CRD v1beta1 and minKubeVersion is informed",
			args: args{
				k8sVersion:     "",
				minKubeVersion: "1.11.3",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantWarning: true,
			warnStrings: []string{"this bundle is " + k8sApiDeprecatedInfo + ". Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
		},
		{
			name: "should not return a warning or error when has minKubeVersion but the k8sVersion informed is <= 1.15",
			args: args{
				k8sVersion:     "1.15",
				minKubeVersion: "1.11.3",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantWarning: true,
			warnStrings: []string{"checking APIs against Kubernetes version : 1.15"},
		},
		{
			name: "should return an error when the k8sVersion is >= 1.22 and has the deprecated API",
			args: args{
				k8sVersion:     "1.22",
				minKubeVersion: "",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantError:   true,
			errStrings:  []string{"this bundle is " + k8sApiDeprecatedInfo + ". Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
			wantWarning: true,
			warnStrings: []string{"checking APIs against Kubernetes version : 1.22"},
		},
		{
			name: "should return an error when the k8sVersion informed is invalid",
			args: args{
				k8sVersion:     "invalid",
				minKubeVersion: "",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantError:  true,
			errStrings: []string{"invalid value informed via the k8s key option : invalid"},
		},
		{
			name: "should return an error when the csv.spec.minKubeVersion informed is invalid",
			args: args{
				minKubeVersion: "invalid",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantError:   true,
			wantWarning: true,
			errStrings:  []string{"unable to use csv.Spec.MinKubeVersion to verify the CRD/Webhook apis because it has an invalid value: invalid"},
			warnStrings: []string{"this bundle is " + k8sApiDeprecatedInfo + ". Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.directory)
			require.NoError(t, err)

			bundle.CSV.Spec.MinKubeVersion = tt.args.minKubeVersion

			errsResult, warnsResult := validateHubDeprecatedAPIS(bundle, tt.args.k8sVersion)

			require.Equal(t, tt.wantWarning, len(warnsResult) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(warnsResult))
				for _, w := range warnsResult {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(errsResult) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(errsResult))
				for _, err := range errsResult {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

func TestValidateHubChannels(t *testing.T) {
	type args struct {
		channels []string
	}
	tests := []struct {
		name     string
		args     args
		wantWarn bool
	}{
		{
			name: "should not return warning when the channel names following the convention",
			args: args{
				channels: []string{"fast", "preview"},
			},
			wantWarn: false,
		},
		{
			name: "should return warning when the channel names are NOT following the convention",
			args: args{
				channels: []string{"mychannel-4.5"},
			},
			wantWarn: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHubChannels(tt.args.channels); (err != nil) != tt.wantWarn {
				t.Errorf("validateHubChannels() error = %v, wantWarn %v", err, tt.wantWarn)
			}
		})
	}
}

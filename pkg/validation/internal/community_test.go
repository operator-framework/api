package internal

import (
	"fmt"
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_communityValidator(t *testing.T) {
	type args struct {
		annotations    map[string]string
		bundleDir      string
		imageIndexPath string
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
			name:      "should work successfully when no deprecated apis are found and has not the annotations or ocp index labels",
			wantError: false,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1",
			},
		},
		{
			name: "should pass when the olm annotation and index label are set with a " +
				"value < 4.9 and has deprecated apis",
			wantError: false,
			args: args{
				bundleDir:      "./testdata/valid_bundle_v1beta1",
				imageIndexPath: "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
		},
		{
			name:      "should fail because is missing the olm.annotation and has deprecated apis",
			wantError: true,
			args: args{
				bundleDir:      "./testdata/valid_bundle_v1beta1",
				imageIndexPath: "./testdata/dockerfile/valid_bundle.Dockerfile",
			},
			errStrings: []string{"Error: Value : (etcdoperator.v0.9.4) csv.Annotations not specified " +
				"olm.maxOpenShiftVersion for an OCP version < 4.9. This annotation is required to " +
				"prevent the user from upgrading their OCP cluster before they have installed a " +
				"version of their operator which is compatible with 4.9. " +
				"This bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22 which are no " +
				"longer supported on 4.9. Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"]) or use the annotation"},
		},
		{
			name: "should fail when the olm annotation is set without the properties for max ocp version and has " +
				"deprecated apis",
			wantError: true,
			args: args{
				bundleDir:      "./testdata/valid_bundle_v1beta1",
				imageIndexPath: "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.invalid", "value": "4.9"}]`),
				},
			},
			errStrings: []string{"Error: Value : (etcdoperator.v0.9.4) csv.Annotations.olm.properties with the key " +
				"`olm.maxOpenShiftVersion` and a value with an OCP version which is < 4.9 is required for any operator " +
				"bundle that is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. " +
				"Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"]) " +
				"or use the annotation"},
		},
		{
			name:      "should fail when the olm annotation is set with a value >= 4.9 and has deprecated apis",
			wantError: true,
			args: args{
				bundleDir:      "./testdata/valid_bundle_v1beta1",
				imageIndexPath: "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.9"}]`),
				},
			},
			errStrings: []string{"Error: Value : (etcdoperator.v0.9.4) csv.Annotations.olm.properties with the key " +
				"and value for olm.maxOpenShiftVersion has the OCP version value 4.9 which is >= of 4.9. " +
				"This bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22." +
				" Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"]) or inform in this property an OCP version which is < 4.9"},
		},
		{
			name:        "should warning because is missing the index-path and has deprecated apis",
			wantWarning: true,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
			warnStrings: []string{"Warning: Value : (etcdoperator.v0.9.4) please, inform the path of its index image " +
				"file via the the optional key values and the key index-path to allow this validator check the labels " +
				"configuration or migrate the API(s) for CRD: " +
				"([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" " +
				"\"etcdrestores.etcd.database.coreos.com\"]). (e.g. index-path=./mypath/bundle.Dockerfile). " +
				"This bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22 "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			if len(tt.args.annotations) > 0 {
				bundle.CSV.Annotations = tt.args.annotations
			}

			results := validateCommunityBundle(bundle, tt.args.imageIndexPath)
			require.Equal(t, tt.wantWarning, len(results.Warnings) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(results.Warnings))
				for _, w := range results.Warnings {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(results.Errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(results.Errors))
				for _, err := range results.Errors {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}

func Test_checkOCPLabelsWithHasDeprecatedAPIs(t *testing.T) {
	type args struct {
		checks    CSVChecks
		indexPath string
	}
	tests := []struct {
		name        string
		args        args
		wantError   bool
		wantWarning bool
	}{
		{
			name: "should pass when has a valid value for the OCP labels",
			args: args{
				indexPath: "./testdata/dockerfile/valid_bundle.Dockerfile",
			},
		},
		{
			name:      "should fail when the OCP label is not found",
			wantError: true,
			args: args{
				indexPath: "./testdata/dockerfile/bundle_without_label.Dockerfile",
			},
		},
		{
			name:      "should fail when the the index path is an invalid path",
			wantError: true,
			args: args{
				indexPath: "./testdata/dockerfile/invalid",
			},
		},
		{
			name:      "should fail when the OCP label index is => 4.9",
			wantError: true,
			args: args{
				indexPath: "./testdata/dockerfile/invalid_bundle_equals_upper.Dockerfile",
			},
		},
		{
			name:      "should fail when the OCP label index range is => 4.9",
			wantError: true,
			args: args{
				indexPath: "./testdata/dockerfile/invalid_bundle_range_upper.Dockerfile",
			},
		},
		{
			name:      "should fail when the OCP label index range is => 4.9 with coma e.g. v4.5,4.6",
			wantError: true,
			args: args{
				indexPath: "./testdata/dockerfile/invalid_bundle_range_upper.Dockerfile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checks := CommunityOperatorChecks{bundle: manifests.Bundle{}, indexImagePath: tt.args.indexPath, errs: []error{}, warns: []error{}}

			checks = checkOCPLabelsWithHasDeprecatedAPIs(checks, "CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])")

			require.Equal(t, tt.wantWarning, len(checks.warns) > 0)
			require.Equal(t, tt.wantError, len(checks.errs) > 0)
		})
	}
}

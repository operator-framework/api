package internal

import (
	"fmt"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
)

func Test_communityValidator(t *testing.T) {
	type args struct {
		annotations   map[string]string
		bundleDir     string
		filePath      string
		ocpLabelRange string
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
				bundleDir: "./testdata/valid_bundle_v1beta1",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
		},
		{
			name: "should pass when the olm annotation and the label in the annotation file is set with a " +
				"value < 4.9 and has deprecated apis",
			wantError: false,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/annotations/annotations.yaml",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
		},
		{
			name: "should pass when the olm annotation and index label are set with a " +
				"value < 4.9 and has deprecated apis and with label flag v4.6-v4.8",
			wantError: false,
			args: args{
				bundleDir:     "./testdata/valid_bundle_v1beta1",
				ocpLabelRange: "v4.6-v4.8",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
		},
		{
			name:      "should fail because is missing the olm.annotation and has deprecated apis",
			wantError: true,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/dockerfile/valid_bundle.Dockerfile",
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
			name:      "should fail when the olm annotation is set with a value >= 4.9 and has deprecated apis",
			wantError: true,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.9"}]`),
				},
			},
			errStrings: []string{"Error: Value : (etcdoperator.v0.9.4) csv.Annotations.olm.properties with the key " +
				"and value for olm.maxOpenShiftVersion has the OCP version value 4.9 which is >= of 4.9. " +
				"This bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22." +
				" Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"]) or inform in this property an OCP version which is < 4.9",
				"Error: Value : (etcdoperator.v0.9.4) the 4.9 annotation to block the cluster upgradedoes not contain in the range v4.6-v4.8 of versions where this solution should be distributed",

			},
		},
		{
			name:        "should warn on patch version in maxOpenShiftVersion",
			wantWarning: true,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8.1"}]`),
				},
			},
			warnStrings: []string{
				"Warning: Value : (etcdoperator.v0.9.4) csv.Annotations.olm.properties has an invalid value. olm.maxOpenShiftVersion must specify only major.minor versions, 4.8.1 will be truncated to 4.8.0",
			},
		},
		{
			name:      "should pass when the maxOpenShiftVersion is semantically equivalent to <major>.<minor>.0",
			wantError: false,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/dockerfile/valid_bundle.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8.0+build"}]`),
				},
			},
		},
		{
			name: "should pass when the olm annotation and index label are set with a " +
				"value =v4.8 and has deprecated apis",
			wantError: false,
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
				filePath:  "./testdata/dockerfile/valid_bundle_4_8.Dockerfile",
				annotations: map[string]string{
					"olm.properties": fmt.Sprintf(`[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]`),
				},
			},
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

			results := validateCommunityBundle(bundle, tt.args.filePath, tt.args.ocpLabelRange)
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
		checks    CommunityOperatorChecks
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
			name:      "should warn when the OCP label is not found",
			wantError: false,
			wantWarning: true,
			args: args{
				indexPath: "./testdata/dockerfile/bundle_without_label.Dockerfile",
			},
		},
		{
			name:      "should fail when the the index path is an invalid path",
			wantError: true,
			wantWarning: true,
			args: args{
				indexPath: "./testdata/dockerfile/invalid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checks := CommunityOperatorChecks{bundle: manifests.Bundle{}, filePath: tt.args.indexPath, errs: []error{}, warns: []error{}}
			checks = getOCPLabel(checks)
			checks = checkOCPLabel(checks)
			require.Equal(t, tt.wantWarning, len(checks.warns) > 0)
			require.Equal(t, tt.wantError, len(checks.errs) > 0)
		})
	}
}

func Test_rangeContainsVersion(t *testing.T) {
	type args struct {
		rangeValue    string
		targetVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "should return true when the label range is <= the targetVersion version",
			wantErr: false,
			want:    true,
			args: args{
				rangeValue:    "=v4.9",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return false when the label range is > than targetVersion version",
			wantErr: false,
			want:    false,
			args: args{
				rangeValue:    "=v4.10",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return false when the label range is > than targetVersion version (v4.10-v4.11) < 4.9",
			wantErr: false,
			want:    false,
			args: args{
				rangeValue:    "v4.10-v4.11",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return invalid syntax",
			wantErr: true,
			want:    false,
			args: args{
				rangeValue:    "’”””v’”4v.vvv’”””8",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return invalid syntax with vv4.vv8v-vvv4vvv.vvv9vvv",
			wantErr: true,
			want:    false,
			args: args{
				rangeValue:    "vv4.vv8v-vvv4vvv.vvv9vvv",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return true when the label range is < than targetVersion version (v4.5-v4.8) < 4.9",
			wantErr: false,
			want:    false,
			args: args{
				rangeValue:    "v4.5-v4.8",
				targetVersion: "4.9",
			},
		},
		{
			name:    "should return true when the label range is > than targetVersion version with comas (v4.5,v4.6) < 4.9",
			wantErr: false,
			want:    true,
			args: args{
				rangeValue:    "v4.5,v4.6",
				targetVersion: "4.9",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rangeContainsVersion(tt.args.rangeValue, tt.args.targetVersion, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("rangeContainsVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rangeContainsVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

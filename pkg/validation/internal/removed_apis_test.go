package internal

import (
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_getDeprecatedAPIs(t *testing.T) {

	// Mock the expected result for ./testdata/valid_bundle_v1beta1
	crdMock := make(map[string][]string)
	crdMock["CRD"] = []string{"etcdbackups.etcd.database.coreos.com", "etcdclusters.etcd.database.coreos.com", "etcdrestores.etcd.database.coreos.com"}

	// Mock the expected result for ./testdata/valid_bundle_with_v1beta1_clusterrole
	otherKindsMock := make(map[string][]string)
	otherKindsMock[ClusterRoleKind] = []string{"memcached-operator-metrics-reader"}
	otherKindsMock[PriorityClassKind] = []string{"super-priority"}
	otherKindsMock[RoleKind] = []string{"memcached-role"}
	otherKindsMock["MutatingWebhookConfiguration"] = []string{"mutating-webhook-configuration"}

	type args struct {
		bundleDir string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return an empty map when no deprecated apis are found",
			args: args{
				bundleDir: "./testdata/valid_bundle_v1",
			},
			want: map[string][]string{},
		},
		{
			name: "should return map with CRDs when this kind of resource is deprecated",
			args: args{
				bundleDir: "./testdata/valid_bundle_v1beta1",
			},
			want: crdMock,
		},
		{
			name: "should return map with others kinds which are deprecated",
			args: args{
				bundleDir: "./testdata/bundle_with_deprecated_resources",
			},
			want: otherKindsMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			if got := getRemovedAPIsOn1_22From(bundle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRemovedAPIsOn1_22From() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDeprecatedAPIS(t *testing.T) {
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
			warnStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. " +
				"Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" " +
				"\"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
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
			wantError: true,
			errStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. " +
				"Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\"" +
				" \"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
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
			errStrings: []string{"unable to use csv.Spec.MinKubeVersion to verify the CRD/Webhook apis because it " +
				"has an invalid value: invalid"},
			warnStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. " +
				"Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" " +
				"\"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.directory)
			require.NoError(t, err)

			bundle.CSV.Spec.MinKubeVersion = tt.args.minKubeVersion

			errsResult, warnsResult := validateDeprecatedAPIS(bundle, tt.args.k8sVersion)

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

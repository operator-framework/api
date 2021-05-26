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

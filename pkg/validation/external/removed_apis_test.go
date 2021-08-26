package external

import (
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_getDeprecated1_25APIs(t *testing.T) {

	// Mock the expected result for ../internal/testdata/valid_bundle_v1beta1
	crdMock := make(map[string][]string)
	crdMock["CRD"] = []string{"etcdbackups.etcd.database.coreos.com", "etcdclusters.etcd.database.coreos.com", "etcdrestores.etcd.database.coreos.com"}

	// Mock the expected result
	otherKindsMock := make(map[string][]string)
	otherKindsMock["PodDisruptionBudget"] = []string{"busybox-pdb"}

	bundleDirPrefix := "../internal/testdata/" // let's reuse the testdata, so we can avoid creating more

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
				bundleDir: bundleDirPrefix + "valid_bundle_v1",
			},
			want: map[string][]string{},
		},
		{
			name: "should return map with others kinds which are deprecated",
			args: args{
				bundleDir: bundleDirPrefix + "bundle_with_deprecated_resources",
			},
			want: otherKindsMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			if got := GetRemovedAPIsOn1_25From(bundle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRemovedAPIsOn1_25From() = %v, want %v", got, tt.want)
			}
		})
	}
}

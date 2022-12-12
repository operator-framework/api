package internal

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
)

func Test_GetRemovedAPIsOn1_22From(t *testing.T) {
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

func Test_GetRemovedAPIsOn1_25From(t *testing.T) {
	mock := make(map[string][]string)
	mock["HorizontalPodAutoscaler"] = []string{"memcached-operator-hpa"}
	mock["PodDisruptionBudget"] = []string{"memcached-operator-policy-manager"}

	warnMock := make(map[string][]string)
	warnMock["cronjobs"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.ClusterPermissions[0].Rules[7]"}
	warnMock["endpointslices"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[3]"}
	warnMock["events"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[2]"}
	warnMock["horizontalpodautoscalers"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[4]"}
	warnMock["poddisruptionbudgets"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[5]"}
	warnMock["podsecuritypolicies"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[5]"}
	warnMock["runtimeclasses"] = []string{"ClusterServiceVersion.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[6]"}

	type args struct {
		bundleDir string
	}
	tests := []struct {
		name     string
		args     args
		errWant  map[string][]string
		warnWant map[string][]string
	}{
		{
			name: "should return an empty map when no deprecated apis are found",
			args: args{
				bundleDir: "./testdata/valid_bundle_v1",
			},
			errWant:  map[string][]string{},
			warnWant: map[string][]string{},
		},
		{
			name: "should fail return the removed APIs in 1.25",
			args: args{
				bundleDir: "./testdata/removed_api_1_25",
			},
			errWant:  mock,
			warnWant: map[string][]string{},
		},
		{
			name: "should return warnings with all deprecated APIs in 1.25",
			args: args{
				bundleDir: "./testdata/deprecated_api_1_25",
			},
			errWant:  mock,
			warnWant: warnMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			errGot, warnGot := getRemovedAPIsOn1_25From(bundle)

			if !reflect.DeepEqual(errGot, tt.errWant) {
				t.Errorf("getRemovedAPIsOn1_25From() = %v, want %v", errGot, tt.errWant)
			}

			if !reflect.DeepEqual(warnGot, tt.warnWant) {
				t.Errorf("getRemovedAPIsOn1_25From() = %v, want %v", warnGot, tt.warnWant)
			}
		})
	}
}

func Test_GetRemovedAPIsOn1_26From(t *testing.T) {
	mock := make(map[string][]string)
	mock["HorizontalPodAutoscaler"] = []string{"memcached-operator-hpa"}

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
			name: "should fail return the removed APIs in 1.26",
			args: args{
				bundleDir: "./testdata/removed_api_1_26",
			},
			want: mock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			if got := getRemovedAPIsOn1_26From(bundle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRemovedAPIsOn1_26From() = %v, want %v", got, tt.want)
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
		},
		{
			name: "should return an error when the k8sVersion is >= 1.25 and found removed APIs on 1.25",
			args: args{
				k8sVersion:     "1.25",
				minKubeVersion: "",
				directory:      "./testdata/removed_api_1_25",
			},
			wantError: true,
			errStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.25. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-25. " +
				"Migrate the API(s) for HorizontalPodAutoscaler: ([\"memcached-operator-hpa\"])," +
				"PodDisruptionBudget: ([\"memcached-operator-policy-manager\"]),"},
		},
		{
			name: "should return a warning if the k8sVersion is empty and found removed APIs on 1.25",
			args: args{
				k8sVersion:     "",
				minKubeVersion: "",
				directory:      "./testdata/removed_api_1_25",
			},
			wantWarning: true,
			warnStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.25. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-25. " +
				"Migrate the API(s) for HorizontalPodAutoscaler: ([\"memcached-operator-hpa\"])," +
				"PodDisruptionBudget: ([\"memcached-operator-policy-manager\"]),"},
		},
		{
			name: "should return an error when the k8sVersion is >= 1.26 and found removed APIs on 1.26",
			args: args{
				k8sVersion:     "1.26",
				minKubeVersion: "",
				directory:      "./testdata/removed_api_1_26",
			},
			wantError: true,
			errStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.26. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-26. " +
				"Migrate the API(s) for HorizontalPodAutoscaler: ([\"memcached-operator-hpa\"])"},
		},
		{
			name: "should return a warning when the k8sVersion is empty and found removed APIs on 1.26",
			args: args{
				k8sVersion:     "",
				minKubeVersion: "",
				directory:      "./testdata/removed_api_1_26",
			},
			wantWarning: true,
			warnStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.26. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-26. " +
				"Migrate the API(s) for HorizontalPodAutoscaler: ([\"memcached-operator-hpa\"])"},
		},
		{
			name: "should return an error when the k8sVersion informed is invalid",
			args: args{
				k8sVersion:     "invalid",
				minKubeVersion: "",
				directory:      "./testdata/valid_bundle_v1beta1",
			},
			wantError:   true,
			errStrings:  []string{"invalid value informed via the k8s key option : invalid"},
			wantWarning: true,
			warnStrings: []string{"this bundle is using APIs which were deprecated and removed in v1.22. " +
				"More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. " +
				"Migrate the API(s) for CRD: ([\"etcdbackups.etcd.database.coreos.com\" " +
				"\"etcdclusters.etcd.database.coreos.com\" \"etcdrestores.etcd.database.coreos.com\"])"},
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
				// testing against sorted strings to address flakiness on the order
				// of APIs listed
				sortedWarnStrings := sortStringSlice(tt.warnStrings)
				for _, w := range warnsResult {
					wString := w.Error()
					require.Contains(t, sortedWarnStrings, sortString(wString))
				}
			}

			require.Equal(t, tt.wantError, len(errsResult) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(errsResult))
				// testing against sorted strings to address flakiness on the order
				// of APIs listed
				sortedErrStrings := sortStringSlice(tt.errStrings)
				for _, err := range errsResult {
					errString := sortString(err.Error())
					require.Contains(t, sortedErrStrings, errString)
				}
			}
		})
	}
}

func sortString(str string) string {
	split := strings.Split(str, "")
	sort.Strings(split)
	return strings.Join(split, "")
}

func sortStringSlice(slice []string) []string {
	var newSlice []string
	for _, str := range slice {
		newSlice = append(newSlice, sortString(str))
	}
	return newSlice
}

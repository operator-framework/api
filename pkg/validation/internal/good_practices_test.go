package internal

import (
	"testing"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
)

func Test_ValidateGoodPractices(t *testing.T) {
	bundleWithDeploymentSpecEmpty, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithDeploymentSpecEmpty.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs = nil

	bundleWithMissingCrdDescription, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithMissingCrdDescription.CSV.Spec.CustomResourceDefinitions.Owned[0].Description = ""

	type args struct {
		bundleDir string
		bundle    *manifests.Bundle
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
			name: "should pass successfully when the resource request is set for " +
				"all containers defined in the bundle",
			args: args{
				bundleDir: "./testdata/valid_bundle",
			},
		},
		{
			name:        "should raise an waring when the resource request is NOT set for any of the containers defined in the bundle",
			wantWarning: true,
			warnStrings: []string{"Warning: Value memcached-operator.v0.0.1: unable to find the resource requests for the container: (kube-rbac-proxy). It is recommended to ensure the resource request for CPU and Memory. Be aware that for some clusters configurations it is required to specify requests or limits for those values. Otherwise, the system or quota may reject Pod creation. More info: https://master.sdk.operatorframework.io/docs/best-practices/managing-resources/",
				"Warning: Value memcached-operator.v0.0.1: unable to find the resource requests for the container: (manager). It is recommended to ensure the resource request for CPU and Memory. Be aware that for some clusters configurations it is required to specify requests or limits for those values. Otherwise, the system or quota may reject Pod creation. More info: https://master.sdk.operatorframework.io/docs/best-practices/managing-resources/"},
			args: args{
				bundleDir: "./testdata/valid_bundle_v1",
			},
		},
		{
			name:      "should fail when the bundle is nil",
			wantError: true,
			args: args{
				bundle: nil,
			},
			errStrings: []string{"Error: : Bundle is nil"},
		},
		{
			name:      "should fail when the bundle csv is nil",
			wantError: true,
			args: args{
				bundle: &manifests.Bundle{CSV: nil, Name: "test"},
			},
			errStrings: []string{"Error: Value test: Bundle csv is nil"},
		},
		{
			name:      "should fail when the csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs is nil",
			wantError: true,
			args: args{
				bundle: bundleWithDeploymentSpecEmpty,
			},
			errStrings: []string{"Error: Value etcdoperator.v0.9.4: unable to find a deployment to install in the CSV"},
		},
		{
			name:        "should raise an warn when the channel does not follows the convention",
			wantWarning: true,
			args: args{
				bundleDir: "./testdata/bundle_with_metadata",
			},
			warnStrings: []string{"Warning: Value memcached-operator.v0.0.1: channel(s) [\"alpha\"] are not following the recommended naming convention: https://olm.operatorframework.io/docs/best-practices/channel-naming"},
		},
		{
			name:        "should raise a warn when a CRD does not have a description",
			wantWarning: true,
			args: args{
				bundle: bundleWithMissingCrdDescription,
			},
			warnStrings: []string{"Warning: Value etcdoperator.v0.9.4: owned CRD \"etcdclusters.etcd.database.coreos.com\" has an empty description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if len(tt.args.bundleDir) > 0 {
				tt.args.bundle, err = manifests.GetBundleFromDir(tt.args.bundleDir)
				require.NoError(t, err)
			}
			results := validateGoodPracticesFrom(tt.args.bundle)
			require.Equal(t, tt.wantWarning, len(results.Warnings) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(results.Warnings))
				for _, w := range results.Warnings {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
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
		name        string
		args        args
		wantWarn    bool
		warnStrings []string
	}{
		{
			name: "should not return warning when the channel names following the convention",
			args: args{
				channels: []string{"fast", "candidate"},
			},
			wantWarn: false,
		},
		{
			name: "should return warning when the channel names are NOT following the convention",
			args: args{
				channels: []string{"mychannel-4.5"},
			},
			wantWarn:    true,
			warnStrings: []string{"channel(s) [\"mychannel-4.5\"] are not following the recommended naming convention: https://olm.operatorframework.io/docs/best-practices/channel-naming"},
		},
		{
			name: "should return warning when has 1 channel NOT following the convention along the others which follows up",
			args: args{
				channels: []string{"alpha", "fast-v2.1", "candidate-v2.2"},
			},
			wantWarn:    true,
			warnStrings: []string{"channel(s) [\"alpha\"] are not following the recommended naming convention: https://olm.operatorframework.io/docs/best-practices/channel-naming"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle, err := manifests.GetBundleFromDir("./testdata/valid_bundle")
			require.NoError(t, err)
			bundle.Channels = tt.args.channels
			err = validateHubChannels(bundle)
			if (err != nil) != tt.wantWarn {
				t.Errorf("validateHubChannels() error = %v, wantWarn %v", err, tt.wantWarn)
			}
			if len(tt.warnStrings) > 0 {
				require.Contains(t, tt.warnStrings, err.Error())
			}
		})
	}
}

func TestValidateRBACForCRDsWith(t *testing.T) {

	bundle, err := manifests.GetBundleFromDir("./testdata/valid_bundle")
	require.NoError(t, err)

	bundleWithPermissions, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithPermissions.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].APIGroups = []string{"apiextensions.k8s.io"}
	bundleWithPermissions.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Resources = []string{"*"}
	bundleWithPermissions.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Verbs = []string{"*"}

	bundleWithPermissionsResource, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithPermissionsResource.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].APIGroups = []string{"apiextensions.k8s.io"}
	bundleWithPermissionsResource.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Resources = []string{"customresourcedefinitions"}
	bundleWithPermissionsResource.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Verbs = []string{"*"}

	bundleWithPermissionsResourceCreate, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithPermissionsResourceCreate.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].APIGroups = []string{"apiextensions.k8s.io"}
	bundleWithPermissionsResourceCreate.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Resources = []string{"customresourcedefinitions"}
	bundleWithPermissionsResourceCreate.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Verbs = []string{"create"}

	bundleWithPermissionsResourcePatch, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithPermissionsResourcePatch.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].APIGroups = []string{"apiextensions.k8s.io"}
	bundleWithPermissionsResourcePatch.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Resources = []string{"customresourcedefinitions"}
	bundleWithPermissionsResourcePatch.CSV.Spec.InstallStrategy.StrategySpec.Permissions[0].Rules[0].Verbs = []string{"patch"}

	type args struct {
		bundleCSV *operatorsv1alpha1.ClusterServiceVersion
	}
	tests := []struct {
		name        string
		args        args
		wantWarn    bool
		warnStrings []string
	}{
		{
			name: "should not return warning when has no permissions",
			args: args{
				bundleCSV: bundle.CSV,
			},
			wantWarn: false,
		},
		{
			name: "should return warning when has permissions for all verbs and resources kind of the apiGroup",
			args: args{
				bundleCSV: bundleWithPermissions.CSV,
			},
			wantWarn:    true,
			warnStrings: []string{"CSV contains permissions to create CRD. An Operator shouldn't deploy or manage other operators (such patterns are known as meta or super operators or include CRDs in its Operands). It's the Operator Lifecycle Manager's job to manage the deployment and lifecycle of operators.  Please, review the design of your solution and if you should not be using Dependency Resolution from OLM instead. More info: https://sdk.operatorframework.io/docs/best-practices/common-recommendation/"},
		},
		{
			name: "should return warning when has permissions for all verbs with the resource specified",
			args: args{
				bundleCSV: bundleWithPermissionsResource.CSV,
			},
			wantWarn:    true,
			warnStrings: []string{"CSV contains permissions to create CRD. An Operator shouldn't deploy or manage other operators (such patterns are known as meta or super operators or include CRDs in its Operands). It's the Operator Lifecycle Manager's job to manage the deployment and lifecycle of operators.  Please, review the design of your solution and if you should not be using Dependency Resolution from OLM instead. More info: https://sdk.operatorframework.io/docs/best-practices/common-recommendation/"},
		},
		{
			name: "should return warning when has permissions to create a CRD",
			args: args{
				bundleCSV: bundleWithPermissionsResourceCreate.CSV,
			},
			wantWarn:    true,
			warnStrings: []string{"CSV contains permissions to create CRD. An Operator shouldn't deploy or manage other operators (such patterns are known as meta or super operators or include CRDs in its Operands). It's the Operator Lifecycle Manager's job to manage the deployment and lifecycle of operators.  Please, review the design of your solution and if you should not be using Dependency Resolution from OLM instead. More info: https://sdk.operatorframework.io/docs/best-practices/common-recommendation/"},
		},
		{
			name: "should return warning when has permissions to create a Patch a CRD",
			args: args{
				bundleCSV: bundleWithPermissionsResourcePatch.CSV,
			},
			wantWarn:    true,
			warnStrings: []string{"CSV contains permissions to create CRD. An Operator shouldn't deploy or manage other operators (such patterns are known as meta or super operators or include CRDs in its Operands). It's the Operator Lifecycle Manager's job to manage the deployment and lifecycle of operators.  Please, review the design of your solution and if you should not be using Dependency Resolution from OLM instead. More info: https://sdk.operatorframework.io/docs/best-practices/common-recommendation/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = validateRBACForCRDsWith(tt.args.bundleCSV)
			if (err != nil) != tt.wantWarn {
				t.Errorf("validateHubChannels() error = %v, wantWarn %v", err, tt.wantWarn)
			}
			if err != nil && len(tt.warnStrings) > 0 {
				require.Contains(t, tt.warnStrings, err.Error())
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
		{
			name: "should not warning for patch releases",
			args: args{bundleName: "memcached-operator.v0.9.2+alpha"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			csv := operatorsv1alpha1.ClusterServiceVersion{}
			csv.Name = tt.args.bundleName
			result := checkBundleName(&csv)

			require.Equal(t, tt.wantWarning, len(result) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(result))
				for _, w := range result {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}
		})
	}
}

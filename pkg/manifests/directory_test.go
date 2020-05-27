package manifests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBundleDir(t *testing.T) {
	bundle, err := GetBundleFromDir("./testdata/valid_bundle")
	require.NoError(t, err)
	require.Equal(t, "etcdoperator.v0.9.4", bundle.Name)
	require.NotNil(t, bundle.CSV)
	require.Equal(t, 3, len(bundle.V1beta1CRDs))
	require.Equal(t, 4, len(bundle.Objects))
}

func TestGetPackage(t *testing.T) {
	pkg, bundles, err := GetManifestsDir("./testdata/valid_package")
	require.NoError(t, err)
	require.NotNil(t, pkg)
	require.Equal(t, "etcd", pkg.PackageName)
	require.Equal(t, 2, len(bundles))
	require.Equal(t, "etcdoperator.v0.9.2", bundles[0].Name)
	require.NotNil(t, bundles[0].CSV)
	require.Equal(t, 3, len(bundles[0].V1beta1CRDs))
	require.Equal(t, "etcdoperator.v0.9.4", bundles[1].Name)
	require.NotNil(t, bundles[1].CSV)
	require.Equal(t, 2, len(bundles[1].V1beta1CRDs))
	require.Equal(t, 1, len(bundles[1].V1CRDs))
}

func TestLoadBundle(t *testing.T) {
	var err error

	_, err = loadBundle("test-operator.v0.0.1", "./testdata/invalid_bundle_with_subdir")
	require.EqualError(t, err, "bundle manifests dir contains directory: testdata/invalid_bundle_with_subdir/foo")
	_, err = loadBundle("test-operator.v0.0.1", "./testdata/invalid_bundle_with_hidden")
	require.EqualError(t, err, "bundle manifests dir has hidden file: testdata/invalid_bundle_with_hidden/.hidden")
}

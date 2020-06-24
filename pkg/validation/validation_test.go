package validation

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation/errors"

	"github.com/stretchr/testify/require"
)

func TestValidateSuccess(t *testing.T) {
	bundle, err := manifests.GetBundleFromDir("./testdata/valid_bundle")
	require.NoError(t, err)

	results := AllValidators.Validate(bundle)
	for _, result := range results {
		require.Equal(t, false, result.HasError())
	}
}

func TestValidate_WithErrors(t *testing.T) {
	bundle, err := manifests.GetBundleFromDir("./testdata/invalid_bundle")
	require.NoError(t, err)

	results := DefaultBundleValidators.Validate(bundle)
	for _, result := range results {
		require.Equal(t, true, result.HasError())
	}
	require.Equal(t, 1, len(results))
}

func TestValidatePackageSuccess(t *testing.T) {
	pkg, bundles, err := manifests.GetManifestsDir("./testdata/valid_package")
	require.NoError(t, err)

	objs := []interface{}{}
	for _, obj := range bundles {
		objs = append(objs, obj)
	}
	objs = append(objs, pkg)

	results := AllValidators.Validate(objs...)
	for _, result := range results {
		require.Equal(t, false, result.HasError())
	}
}

func TestValidatePackageError(t *testing.T) {
	pkg, _, err := manifests.GetManifestsDir("./testdata/invalid_package")
	require.NoError(t, err)

	objs := []interface{}{pkg}

	results := AllValidators.Validate(objs...)
	require.Equal(t, 1, len(results))
	require.True(t, results[0].HasError())
	pkgErrs := results[0].Errors
	require.Equal(t, 1, len(pkgErrs))
	require.Equal(t, errors.ErrInvalidPackageManifest("packageName empty", pkg.PackageName), pkgErrs[0])
}

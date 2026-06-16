package manifests

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBundleLoaderLoadBundleWalkFunc_PropagatesWalkError(t *testing.T) {
	loader := NewBundleLoader(t.TempDir())
	walkErr := errors.New("permission denied")

	err := loader.LoadBundleWalkFunc("some/path", nil, walkErr)
	require.ErrorIs(t, err, walkErr)
}

func TestPackageManifestLoaderLoadPackagesWalkFunc_PropagatesWalkError(t *testing.T) {
	loader := NewPackageManifestLoader(t.TempDir())
	walkErr := errors.New("permission denied")

	err := loader.LoadPackagesWalkFunc("some/path", nil, walkErr)
	require.ErrorIs(t, err, walkErr)
}

func TestPackageManifestLoaderLoadBundleWalkFunc_PropagatesWalkError(t *testing.T) {
	loader := NewPackageManifestLoader(t.TempDir())
	walkErr := errors.New("permission denied")

	err := loader.LoadBundleWalkFunc("some/path", nil, walkErr)
	require.ErrorIs(t, err, walkErr)
}

func TestCollectWalkErrs_CollectsWalkErrors(t *testing.T) {
	walkErr := errors.New("stat failed")
	loader := NewBundleLoader(t.TempDir())

	var errs []error
	wrapped := collectWalkErrs(loader.LoadBundleWalkFunc, &errs)

	// Calling with a non-nil walk error should collect it and not panic.
	err := wrapped("some/path", nil, walkErr)
	require.NoError(t, err, "collectWalkErrs should swallow non-SkipDir errors so Walk continues")
	require.Len(t, errs, 1)
	require.ErrorIs(t, errs[0], walkErr)
}

func TestBundleLoaderLoadBundleWalkFunc_NilFileInfoWithoutError(t *testing.T) {
	loader := NewBundleLoader(t.TempDir())

	err := loader.LoadBundleWalkFunc("some/path", nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid file")
}

func TestLoadBundle_NonexistentDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")

	loader := NewBundleLoader(dir)
	err := loader.LoadBundle()
	require.Error(t, err)
}

func TestLoadPackage_NonexistentDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")

	loader := NewPackageManifestLoader(dir)
	err := loader.LoadPackage()
	require.Error(t, err)
}

func TestLoadBundle_InaccessibleSubpath(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("test requires non-root: root bypasses directory permission checks")
	}
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	require.NoError(t, os.Mkdir(sub, 0o000))
	t.Cleanup(func() { os.Chmod(sub, 0o755) })

	loader := NewBundleLoader(dir)
	err := loader.LoadBundle()
	require.Error(t, err)
}

func TestLoadPackage_InaccessibleSubpath(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("test requires non-root: root bypasses directory permission checks")
	}
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	require.NoError(t, os.Mkdir(sub, 0o000))
	t.Cleanup(func() { os.Chmod(sub, 0o755) })

	loader := NewPackageManifestLoader(dir)
	err := loader.LoadPackage()
	require.Error(t, err)
}

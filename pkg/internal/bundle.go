package manifests

import (
	"encoding/json"

	manifests "github.com/operator-framework/api/pkg/registry/manifests"

	"github.com/blang/semver"
	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/operator-framework/operator-registry/pkg/sqlite"
	"github.com/pkg/errors"
)

// manifestsLoad loads a manifests directory from disk.
type manifestsLoad struct {
	dir     string
	pkg     registry.PackageManifest
	bundles map[string]*registry.Bundle
}

// Ensure manifestsLoad implements registry.Load.
var _ registry.Load = &manifestsLoad{}

// populate uses operator-registry's sqlite.NewSQLLoaderForDirectory to load
// l.dir's manifests. Note that this method does not call any functions that
// use SQL drivers.
func (l *manifestsLoad) populate() error {
	loader := sqlite.NewSQLLoaderForDirectory(l, l.dir)
	if err := loader.Populate(); err != nil {
		return errors.Wrapf(err, "error getting bundles from manifests dir %q", l.dir)
	}
	return nil
}

// AddOperatorBundle adds a bundle to l.
func (l *manifestsLoad) AddOperatorBundle(bundle *registry.Bundle) error {
	csvRaw, err := bundle.ClusterServiceVersion()
	if err != nil {
		return errors.Wrap(err, "error getting bundle CSV")
	}
	csvSpec := operatorsv1alpha1.ClusterServiceVersionSpec{}
	if err := json.Unmarshal(csvRaw.Spec, &csvSpec); err != nil {
		return errors.Wrap(err, "error unmarshaling CSV spec")
	}
	bundle.Name = csvSpec.Version.String()
	l.bundles[csvSpec.Version.String()] = bundle
	return nil
}

// AddOperatorBundle adds the package manifest to l.
func (l *manifestsLoad) AddPackageChannels(pkg registry.PackageManifest) error {
	l.pkg = pkg
	return nil
}

// AddBundlePackageChannels is a no-op to implement the registry.Load interface.
func (l *manifestsLoad) AddBundlePackageChannels(manifest registry.PackageManifest, bundle registry.Bundle) error {
	return nil
}

// RmPackageName is a no-op to implement the registry.Load interface.
func (l *manifestsLoad) RmPackageName(packageName string) error {
	return nil
}

// ClearNonDefaultBundles is a no-op to implement the registry.Load interface.
func (l *manifestsLoad) ClearNonDefaultBundles(packageName string) error {
	return nil
}

// ManifestsStorer knows how to query for an operator's package manifest and
// related bundles.
type ManifestsStorer interface {
	// GetPackageManifest returns the ManifestsStorer's registry.PackageManifest.
	// The returned object is assumed to be valid.
	GetPackageManifest() manifests.PackageManifest
	// GetBundles returns the ManifestsStorer's set of Bundles. These bundles
	// are unique by CSV version, since only one operator type should exist
	// in one manifests dir.
	// The returned objects are assumed to be valid.
	GetBundles() []*manifests.Bundle
	// GetBundleForVersion returns the ManifestsStorer's Bundle for a given
	// version string. An error should be returned if the passed version
	// does not exist in the store.
	// The returned object is assumed to be valid.
	GetBundleForVersion(string) (*manifests.Bundle, error)
}

// manifestsStore implements ManifestsStorer
type manifestsStore struct {
	pkg     manifests.PackageManifest
	bundles map[string]*manifests.Bundle
}

// ManifestsStoreForDir populates a ManifestsStorer from the metadata in dir.
// Each bundle and the package manifest are statically validated, and will
// return an error if any are not valid.
func ManifestsStoreForDir(dir string) (ManifestsStorer, error) {
	load := &manifestsLoad{
		dir:     dir,
		bundles: map[string]*registry.Bundle{},
	}
	if err := load.populate(); err != nil {
		return nil, err
	}
	// TODO(estroz): remove when operator-registry migrates to api types.
	pkg, bundles := convertRegistryToAPITypes(load.pkg, load.bundles)
	return &manifestsStore{
		pkg:     pkg,
		bundles: bundles,
	}, nil
}

// TODO(estroz): remove when operator-registry migrates to api types.
func convertRegistryToAPITypes(pkgR registry.PackageManifest, bundlesR map[string]*registry.Bundle) (manifests.PackageManifest, map[string]*manifests.Bundle) {
	pkgA := manifests.PackageManifest{
		PackageName:        pkgR.PackageName,
		DefaultChannelName: pkgR.DefaultChannelName,
	}
	for _, channel := range pkgR.Channels {
		pkgA.Channels = append(pkgA.Channels, manifests.PackageChannel{
			Name:           channel.Name,
			CurrentCSVName: channel.CurrentCSVName,
		})
	}
	bundlesA := map[string]*manifests.Bundle{}
	for key, bundle := range bundlesR {
		b := manifests.Bundle{
			Name:        bundle.Name,
			Package:     bundle.Package,
			Channel:     bundle.Channel,
			BundleImage: bundle.BundleImage,
		}
		for _, obj := range bundle.Objects {
			b.Add(obj.DeepCopy())
		}
		bundlesA[key] = &b
	}
	return pkgA, bundlesA
}

func (s manifestsStore) GetPackageManifest() manifests.PackageManifest {
	return s.pkg
}

func (s manifestsStore) GetBundles() (bundles []*manifests.Bundle) {
	for _, bundle := range s.bundles {
		bundles = append(bundles, bundle)
	}
	return bundles
}

func (s manifestsStore) GetBundleForVersion(version string) (*manifests.Bundle, error) {
	if _, err := semver.Parse(version); err != nil {
		return nil, errors.Wrapf(err, "error getting bundle for version %q", version)
	}
	bundle, ok := s.bundles[version]
	if !ok {
		return nil, errors.Errorf("bundle for version %q does not exist", version)
	}
	return bundle, nil
}

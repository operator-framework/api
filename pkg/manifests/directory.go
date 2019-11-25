package manifests

import (
	"fmt"

	internal "github.com/operator-framework/api/pkg/internal"
	manifests "github.com/operator-framework/api/pkg/registry/manifests"
	"github.com/operator-framework/api/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"
)

// GetManifestsDir parses all bundles and a package manifest from dir, which
// are returned if found along with any errors or warnings encountered while
// parsing/validating found manifests.
func GetManifestsDir(dir string) (manifests.PackageManifest, []*manifests.Bundle, []errors.ManifestResult) {
	store, err := internal.ManifestsStoreForDir(dir)
	if err != nil {
		result := errors.ManifestResult{}
		result.Add(errors.ErrInvalidParse(fmt.Sprintf("parse manifests from %q", dir), err))
		return manifests.PackageManifest{}, nil, []errors.ManifestResult{result}
	}
	pkg := store.GetPackageManifest()
	bundles := store.GetBundles()
	objs := []interface{}{}
	for _, obj := range bundles {
		objs = append(objs, obj)
	}
	results := validation.AllValidators.Validate(objs...)
	return pkg, bundles, results
}

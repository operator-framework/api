// Package validation provides default Validator's that can be run with a list
// of arbitrary objects. The defaults exposed here consist of all Validator's
// implemented by this validation library.
//
// Each default Validator runs an independent set of validation functions on
// a set of objects. To run all implemented Validator's, use DefaultValidators().
// The Validator will not be run on objects not of the appropriate type.

package validation

import (
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
	"github.com/operator-framework/api/pkg/validation/internal"
)

// DefaultPackageManifestValidators returns a package manifest Validator that
// can be run directly with Apply(objs...). Optionally, any additional
// Validator's can be added to the returned Validators set.
func DefaultPackageManifestValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals, internal.PackageManifestValidator{})
}

// DefaultClusterServiceVersionValidators returns a ClusterServiceVersion
// Validator that can be run directly with Apply(objs...). Optionally, any
// additional Validator's can be added to the returned Validators set.
func DefaultClusterServiceVersionValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals, internal.CSVValidator{})
}

// DefaultCustomResourceDefinitionValidators returns a CustomResourceDefinition
// Validator that can be run directly with Apply(objs...). Optionally, any
// additional Validator's can be added to the returned Validators set.
func DefaultCustomResourceDefinitionValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals, internal.CRDValidator{})
}

// DefaultBundleValidators returns a bundle Validator that can be run directly
// with Apply(objs...). Optionally, any additional Validator's can be added to
// the returned Validators set.
func DefaultBundleValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals, internal.BundleValidator{})
}

// DefaultManifestsValidators returns a manifests Validator that can be run
// directly with Apply(objs...). Optionally, any additional Validator's can be
// added to the returned Validators set.
func DefaultManifestsValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals, internal.ManifestsValidator{})
}

// DefaultValidators returns all default Validator's, which can be run directly
// with Apply(objs...). Optionally, any additional Validator's can be added to
// the returned Validators set.
func DefaultValidators(vals ...interfaces.Validator) interfaces.Validators {
	return append(vals,
		internal.PackageManifestValidator{},
		internal.CSVValidator{},
		internal.CRDValidator{},
		internal.BundleValidator{},
		internal.ManifestsValidator{},
	)
}

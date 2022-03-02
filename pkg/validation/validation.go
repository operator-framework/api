// Package validation provides default Validator's that can be run with a list
// of arbitrary objects. The defaults exposed here consist of all Validator's
// implemented by this validation library.
//
// Each default Validator runs an independent set of validation functions on
// a set of objects. To run all implemented Validator's, use AllValidators.
// The Validator will not be run on objects of an inappropriate type.

package validation

import (
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
	"github.com/operator-framework/api/pkg/validation/internal"
)

// PackageManifestValidator implements Validator to validate package manifests.
var PackageManifestValidator = internal.PackageManifestValidator

// ClusterServiceVersionValidator implements Validator to validate
// ClusterServiceVersions.
var ClusterServiceVersionValidator = internal.CSVValidator

// CustomResourceDefinitionValidator implements Validator to validate
// CustomResourceDefinitions.
var CustomResourceDefinitionValidator = internal.CRDValidator

// BundleValidator implements Validator to validate Bundles.
//
// This check will verify if the Bundle spec is valid by checking:
//
// - for duplicate keys in the bundle, which may occur if a v1 and v1beta1 CRD of the same GVK appear.
//
// - if owned keys must matches with a CRD in bundle
//
// - if the bundle has APIs(CRDs) which are not defined in the CSV
//
// - if the bundle size compressed is < ~1MB
//
// NOTE: The bundle size test will raise an error if the size is bigger than the max allowed
// and warnings when:
// a) the api is unable to check the bundle size because we are running a check without load the bundle
//
// b) the api could identify that the bundle size is close to the limit (bigger than 85%)
//
// c) [Deprecated and planned to be removed at 2023 -  The API will start growing to encompass validation for all past history] - if the bundle size uncompressed < ~1MB and it cannot work on clusters which uses OLM versions < 1.17.5
//
var BundleValidator = internal.BundleValidator

// OperatorHubValidator implements Validator to validate bundle objects
// for OperatorHub.io requirements.
var OperatorHubValidator = internal.OperatorHubValidator

// Object Validator validates various custom objects in the bundle like PDBs and SCCs.
// Object validation is optional and not a default-level validation.
var ObjectValidator = internal.ObjectValidator

// OperatorGroupValidator implements Validator to validate OperatorGroup manifests
var OperatorGroupValidator = internal.OperatorGroupValidator

// CommunityOperatorValidator implements Validator to validate bundle objects
// for the Community Operator requirements.
//
// Deprecated - The checks made for this validator were moved to the external one:
// https://github.com/redhat-openshift-ecosystem/ocp-olm-catalog-validator.
// Please no longer use this check it will be removed in the next releases.
var CommunityOperatorValidator = internal.CommunityOperatorValidator

// AlphaDeprecatedAPIsValidator implements Validator to validate bundle objects
// for API deprecation requirements.
//
// Note that this validator looks at the manifests. If any removed APIs for the mapped k8s versions are found,
// it raises a warning.
//
// This validator only raises an error when the deprecated API found is removed in the specified k8s
// version informed via the optional key `k8s-version`.
var AlphaDeprecatedAPIsValidator = internal.AlphaDeprecatedAPIsValidator

// GoodPracticesValidator implements Validator to validate the criteria defined as good practices
var GoodPracticesValidator = internal.GoodPracticesValidator

// AllValidators implements Validator to validate all Operator manifest types.
var AllValidators = interfaces.Validators{
	PackageManifestValidator,
	ClusterServiceVersionValidator,
	CustomResourceDefinitionValidator,
	BundleValidator,
	OperatorHubValidator,
	ObjectValidator,
	OperatorGroupValidator,
	CommunityOperatorValidator,
	AlphaDeprecatedAPIsValidator,
	GoodPracticesValidator,
}

var DefaultBundleValidators = interfaces.Validators{
	ClusterServiceVersionValidator,
	CustomResourceDefinitionValidator,
	BundleValidator,
}

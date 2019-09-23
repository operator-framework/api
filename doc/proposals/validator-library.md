# Title

Status: `Draft`

Version: `alpha`

Implementation Owner: @estroz

#### Table of Contents
- [Summary](#summary)
- [Motivation](#motivation)
    - [Goals](#goals)
    - [Non-goals](#non-goals)
- [Proposal](#proposal)
    - [Design Overview](#design-overview)
    - [Implementation details](#implementation-details)
    - [Risks and mitigations](#risks-and-mitigations)
- [Observations and open questions](#observations-and-open-questions)

## Summary

Manifests can be difficult to write correctly because the objects they create often are not statically checked for correctness. Most Operator developers do not realize their manifest is invalid until `kubectl apply -f` informs them them. Static validation can prevent many of these errors from being discovered at runtime by instead bringing them to the user's attention during manifest development.

## Motivation

A general static validation library does not exist for Kubernetes manifests. Instead, validating manifest keys and data is often done via `kubectl apply -f`. This works but does not fully validate manifest data (the object may create but error during installation) and is slow. A validation library with sane defaults for commonly used object schemas and an extensible interface is needed to speed up development workflows and centralize validation knowledge.

### Goals

- Static validation of Kubernetes manifests on-disk.
- A useful set of default validators.
- Returns useful errors containing an object ID, error type, and fields/values related to the error.
- Validation interface is extensible, such that anyone can implement their own and add it dynamically to the default set of validators.

### Non-goals

- Implement validators for all Kubernetes objects.
- Fully replace `kubectl apply -f`.

## Proposal

### Design Overview

#### Interface composition

A `Validator` interface will be exposed that requires a function be implemented that takes a set of objects and returns closure functions around objects with the desired types. These closures can then be collected and run, returning a `ManifestResult` type containing errors and warnings related to incorrectly formatted manifest fields/data.

One closure function around one object is a unit of validation. This discretization allows developers to break up validation logic into a comprehensible manner and turn validators on/off by modifying single lines of code. Library users can compose multiple validator units into one `Validators` object, and run them all at once on one or more objects. Running validators can be done lazily (collect validator closures and run them later) or actively (invoke the closure at the point of its return).

#### Errors

An `errors` package will be provided that exposes custom result types for validation. The `ManifestResult` type will be returned by each closure function, containing sets of `Error`'s and `Warning`'s. The underlying `Error` type for both sets will be similar to [`field.Error`][crd-validation-error]. These results will help end users discover what needs fixing in their manifests.

#### Default `Validator`'s

A default set of `Validator` implementers will be available:

- `ClusterServiceVersion`: validate `ClusterServiceVersion` objects.
- `CustomResourceDefinition`: validate `CustomResourceDefinition` objects.
- Package manifest: validate package manifest objects.
- Bundle: validate all objects in a bundle directory (via [`*registry.Bundle`][registry-bundle-type]) collectively, verifying relationships between objects.
- Manifests: validate all objects in a manifests directory collectively, verifying relationships between objects.

Each will be exported as a `Default{Type}Validators` function containing a set of validators disjointed from other default validators. These functions should take a variadic `...Validator` parameter, such that a user can add custom validators to the list of defaults provided.

### Implementation details

`Validator` interface:

```go
package validator

import (
	"github.com/operator-framework/api/pkg/validation/errors"
)

// ValidatorFunc returns a ManifestResult containing errors and warnings from
// validating some object. These are typically closures or methods.
type ValidatorFunc func() errors.ManifestResult

// ValidationFuncs is a set of validation functions.
type ValidatorFuncs []ValidatorFunc

// Validator is an interface for validating arbitrary objects. A Validator
// returns a set of functions that each validate some underlying object.
// This allows the implementer to easily break down each validation step
// into discrete units, and perform either lazy or active validation.
type Validator interface {
	// GetFuncs takes a list of arbitrary objects and returns a set of functions
	// that each validate some object from the list.
	GetFuncs(...interface{}) ValidatorFuncs
}

// Validators is a set of Validator's that can be run via Apply.
type Validators []Validator

// Apply collects validator functions from each Validator in vals by collecting
// the appropriate functions for each obj in objs, then invokes them,
// collecting the results. Use Apply in code that:
// - Uses more than one Validator in one call.
// - Want active validation.
func (vals Validators) Apply(objs ...interface{}) (results []errors.ManifestResult) {
	for _, val := range vals {
		for _, validate := range val.GetFuncs(objs...) {
			results = append(results, validate())
		}
	}
	return results
}
```

A prototypical implementation for `CustomResourceDefinition`'s:

```go
package validation

import (
	"github.com/operator-framework/api/pkg/validation/errors"
	"github.com/operator-framework/api/pkg/validation/validator"

	apiextv1beta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type crdValidator struct{}

func DefaultCustomResourceDefinitionValidators(vals ...validator.Validator) validator.Validators {
	return append(vals, crdValidator{})
}

func (f crdValidator) GetFuncs(objs ...interface{}) (funcs validator.ValidatorFuncs) {
	for _, obj := range objs {
		switch v := obj.(type) {
		case *apiextv1beta.CustomResourceDefinition:
			funcs = append(funcs, func() errors.ManifestResult {
				return validateCRD(v)
			})
		}
	}
	return funcs
}

func validateCRD(crd interface{}) errors.ManifestResult {
	...
}
```

A more complex implementation, involving a validator for multiple object types:

```go
package validation

import (
	"github.com/operator-framework/api/pkg/validation/errors"
	"github.com/operator-framework/api/pkg/validation/validator"

	"github.com/operator-framework/operator-registry/pkg/registry"
)

type manifestsValidator struct{}

func DefaultManifestsValidators(vals ...validator.Validator) validator.Validators {
	return append(vals, manifestsValidator{})
}

func (f manifestsValidator) GetFuncs(objs ...interface{}) (funcs validator.ValidatorFuncs) {
	var pkg *registry.PackageManifest
	bundles := []*registry.Bundle{}
	for _, obj := range objs {
		switch v := obj.(type) {
		case *registry.PackageManifest:
			if pkg == nil {
				pkg = v
			}
		case *registry.Bundle:
			bundles = append(bundles, v)
		}
	}
	if pkg != nil && len(bundles) > 0 {
		funcs = append(funcs, func() errors.ManifestResult {
			return validateManifests(pkg, bundles)
		})
	}
	return funcs
}

// Use data from pkg to ensure bundles are valid.
func validateManifests(pkg *registry.PackageManifest, bundles []*registry.Bundle) (result errors.ManifestResult) {
	...
}
```

### Risks and mitigations

## Observations and open questions

- Is this interface setup too heavy-duty? Do we want something simpler?
- Does this interface fulfill the outlined goals?


[crd-validation-error]:https://godoc.org/k8s.io/apimachinery/pkg/util/validation/field#Error
[registry-bundle]:https://github.com/operator-framework/operator-registry/blob/9d997b8/pkg/registry/bundle.go#L32

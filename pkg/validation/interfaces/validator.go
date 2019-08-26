package validator

import (
	"github.com/operator-framework/api/pkg/validation/errors"
)

// Validator is an interface for validating arbitrary objects.
type Validator interface {
	// Validate takes a list of arbitrary objects and returns a slice of results,
	// one for each object validated.
	Validate(...interface{}) []errors.ManifestResult
}

// Validators is a set of Validator's that can be run via Apply.
type Validators []Validator

// ApplyParallel invokes each Validator in vals, and collects and returns
// the results.
func (vals Validators) Apply(objs ...interface{}) (results []errors.ManifestResult) {
	for _, validator := range vals {
		results = append(results, validator.Validate(objs...)...)
	}
	return results
}

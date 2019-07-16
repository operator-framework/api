package validator

import "fmt"

// ManifestResult represents verification result for each of the yaml files
// from the manifest bundle.
type ManifestResult struct {
	// Name is some piece of information identifying the manifest. This should
	// usually be set to object.GetName().
	Name string
	// Errors pertain to issues with the manifest that must be corrected.
	Errors []MissingTypeError
	// Warnings pertain to issues with the manifest that are optional to correct.
	Warnings []MissingTypeError
}

// MissingTypeError represents a warning or an error in a yaml file.
type MissingTypeError struct {
	// Err is the underlying error caused by a missing type, if any.
	Err error
	// TypeName is the syntactical name of missing data, ex. Struct, Field.
	TypeName string
	// Path is the dot-hierarchical YAML path of the missing data.
	Path string
	// IsMandatory determines whether the missing data should generate a
	// warning (false, the default) or error (true).
	IsMandatory bool
}

// MissingTypeError strut implements the Error interface to define custom error formatting.
func (err MissingTypeError) Error() string {
	if err.IsMandatory {
		return fmt.Sprintf("Error: Mandatory %s Missing (%s)", err.TypeName, err.Path)
	} else {
		return fmt.Sprintf("Warning: Optional %s Missing (%s)", err.TypeName, err.Path)
	}
}

// ValidatorSet contains a set of Validators to be executed sequentially.
// TODO: add configurable logger.
type ValidatorSet struct {
	validators []Validator
}

// NewValidatorSet creates a ValidatorSet containing vs.
func NewValidatorSet(vs ...Validator) *ValidatorSet {
	set := &ValidatorSet{}
	set.AddValidators(vs...)
	return set
}

// AddValidators adds each unique Validator in vs to the receiver.
func (set *ValidatorSet) AddValidators(vs ...Validator) {
	seenNames := map[string]struct{}{}
	for _, v := range vs {
		if _, seen := seenNames[v.Name()]; !seen {
			set.validators = append(set.validators, v)
			seenNames[v.Name()] = struct{}{}
		}
	}
}

// ValidateAll runs each Validator in the receiver and returns all results.
func (set ValidatorSet) ValidateAll() (allResults []ManifestResult) {
	for _, v := range set.validators {
		results := v.Validate()
		allResults = append(allResults, results...)
	}
	return allResults
}

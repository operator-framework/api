package validator

// Validator is an interface for implementing a validator of a single
// Kubernetes object type. Ideally each Validator will check one aspect of
// an object, or perform several steps that have a common theme or goal.
type Validator interface {
	// Validate should run validation logic on an arbitrary object, and return
	// a one ManifestResult for each object that did not pass validation.
	// TODO: use pointers
	Validate() []ManifestResult
	// AddObjects adds objects to the Validator. Each object will be validated
	// when Validate() is called.
	AddObjects(...interface{}) Error
	// Name should return a succinct name for this validator.
	Name() string
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

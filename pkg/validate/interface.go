package validate

// Validator is an interface for implementing a validator of a single
// Kubernetes object type. Ideally each Validator will check one aspect of
// an object, or perform several steps that have a common theme or goal.
type Validator interface {
	// Validate should run validation logic on an arbitrary object, and return
	// an error if the object does not pass validation.
	Validate() error
	// AddObjects adds objects to the Validator. Each object will be validated
	// when Validate() is called.
	AddObjects(...interface{}) error
	// Name should return a succinct name for this validator.
	Name() string
}

type ValidatorSet struct {
	validators []Validator
}

func NewValidatorSet(vs ...Validator) *ValidatorSet {
	set := &ValidatorSet{}
	set.AddValidators(vs...)
	return set
}

func (set *ValidatorSet) AddValidators(vs ...Validator) {
	seenNames := map[string]struct{}{}
	for _, v := range vs {
		if _, seen := seenNames[v.Name()]; !seen {
			set.validators = append(set.validators, v)
			seenNames[v.Name()] = struct{}{}
		}
	}
}

func (set ValidatorSet) ValidateAll() (errs []error) {
	for _, v := range set.validators {
		if err := v.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

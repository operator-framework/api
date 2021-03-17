package schema

import (
	"github.com/sirupsen/logrus"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
)

type ConfigValidator interface {
	Validate(b []byte, key string) error
}

// NewConfigValidator is a constructor that returns a ConfigValidator
func NewConfigValidator(instance *build.Instance, logger *logrus.Entry) ConfigValidator {
	return configValidator{
		runtime:  &cue.Runtime{},
		instance: instance,
		logger:   logger,
	}
}

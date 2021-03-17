package schema

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/encoding/json"
)

const (
	olmBundle          = "olmbundle"
	olmPackage         = "olmpackage"
	olmChannel         = "olmchannel"
	olmSkips           = "olmskips"
	olmSkipRange       = "olmskipRange"
	olmGVKProvided     = "olmgvkprovided"
	olmGVKRequired     = "olmgvkrequired"
	olmPackageProperty = "packageproperty"
)

type configValidator struct {
	runtime  *cue.Runtime
	instance *build.Instance
	logger   *logrus.Entry
}

func (c configValidator) Validate(b []byte, key string) error {
	inst, err := c.runtime.Build(c.instance)
	if err != nil {
		return err
	}

	v := inst.LookupDef(key)
	if !v.Exists() {
		err := fmt.Errorf("Unable to find the definition %s in schema", key)
		c.logger.WithError(err).Debugf(key)
		return err
	}

	err = json.Validate(b, v)
	if err != nil {
		c.logger.WithError(err).Debugf("Validation error")
		return err
	}

	return nil
}

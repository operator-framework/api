package validate

import (
	"encoding/json"
	"fmt"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"
)

type PackageValidator struct {
	fileName string
	pkgs     []registry.PackageManifest
}

var _ validator.Validator = &PackageValidator{}

// packageManifest is an alias of `registry.PackageManifest` used
// to define new methods on the non-local type.
type packageManifest registry.PackageManifest

func (p packageManifest) getPkgName() string {
	return p.PackageName
}

func (v *PackageValidator) Validate() (results []validator.ManifestResult) {
	for _, pkg := range v.pkgs {
		result := pkgInspect(pkg)
		if result.Name == "" {
			result.Name = packageManifest(pkg).getPkgName()
		}
		results = append(results, result)
	}
	return results
}

func (v *PackageValidator) AddObjects(objs ...interface{}) validator.Error {
	for _, o := range objs {
		switch t := o.(type) {
		case registry.PackageManifest:
			v.pkgs = append(v.pkgs, t)
		case *registry.PackageManifest:
			v.pkgs = append(v.pkgs, *t)
		}
	}
	return validator.Error{}
}

func (v PackageValidator) Name() string {
	return "Package Validator"
}

func (v PackageValidator) FileName() string {
	return v.fileName
}

func (v PackageValidator) Unmarshal(rawYaml []byte) (interface{}, error) {
	var pkg registry.PackageManifest

	rawJson, err := yaml.YAMLToJSON(rawYaml)
	if err != nil {
		return registry.PackageManifest{}, fmt.Errorf("error parsing raw YAML to Json: %s", err)
	}
	if err := json.Unmarshal(rawJson, &pkg); err != nil {
		return registry.PackageManifest{}, fmt.Errorf("error parsing package type (JSON) : %s", err)
	}
	return pkg, nil
}

func pkgInspect(pkg registry.PackageManifest) (manifestResult validator.ManifestResult) {
	manifestResult = validator.ManifestResult{}
	present, manifestResult := isDefaultPresent(pkg, manifestResult)
	if !present {
		manifestResult.Errors = append(manifestResult.Errors, validator.InvalidDefaultChannel(fmt.Sprintf("Error: default channel %s not found in the list of declared channels", pkg.DefaultChannelName), pkg.DefaultChannelName))
	}
	return
}

func isDefaultPresent(pkg registry.PackageManifest, manifestResult validator.ManifestResult) (bool, validator.ManifestResult) {
	present := false
	for _, channel := range pkg.Channels {
		if pkg.DefaultChannelName == "" {
			manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidDefaultChannel(fmt.Sprintf("Warning: default channel not found in %s package manifest", pkg.PackageName), pkg.PackageName))
			return true, manifestResult
		} else if pkg.DefaultChannelName == channel.Name {
			present = true
		}
	}
	return present, manifestResult
}

package types

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
)

const (
	GVKType        = "olm.gvk"
	PackageType    = "olm.package"
	DeprecatedType = "olm.deprecated"
	LabelType      = "olm.label"
	PropertyKey    = "olm.properties"
)

// DependenciesFile holds dependency information about a bundle
type DependenciesFile struct {
	// Dependencies is a list of dependencies for a given bundle
	Dependencies []Dependency `json:"dependencies" yaml:"dependencies"`
}

// Dependency specifies a single constraint that can be satisfied by a property on another bundle..
type Dependency struct {
	// The type of dependency. This field is required.
	Type string `json:"type" yaml:"type"`

	// The serialized value of the dependency
	Value json.RawMessage `json:"value" yaml:"value"`
}

type GVKDependency struct {
	// The group of GVK based dependency
	Group string `json:"group" yaml:"group"`

	// The kind of GVK based dependency
	Kind string `json:"kind" yaml:"kind"`

	// The version of GVK based dependency
	Version string `json:"version" yaml:"version"`
}

type PackageDependency struct {
	// The name of dependency such as 'etcd'
	PackageName string `json:"packageName" yaml:"packageName"`

	// The version range of dependency in semver range format
	Version string `json:"version" yaml:"version"`
}

type LabelDependency struct {
	// The Label name of dependency
	Label string `json:"label" yaml:"label"`
}


// Validate will validate GVK dependency type and return error(s)
func (gd *GVKDependency) Validate() []error {
	errs := []error{}
	if gd.Group == "" {
		errs = append(errs, fmt.Errorf("API Group is empty"))
	}
	if gd.Version == "" {
		errs = append(errs, fmt.Errorf("API Version is empty"))
	}
	if gd.Kind == "" {
		errs = append(errs, fmt.Errorf("API Kind is empty"))
	}
	return errs
}

// Validate will validate Label dependency type and return error(s)
func (ld *LabelDependency) Validate() []error {
	errs := []error{}
	if *ld == (LabelDependency{}) {
		errs = append(errs, fmt.Errorf("Label information is missing"))
	}
	return errs
}

// Validate will validate package dependency type and return error(s)
func (pd *PackageDependency) Validate() []error {
	errs := []error{}
	if pd.PackageName == "" {
		errs = append(errs, fmt.Errorf("Package name is empty"))
	}
	if pd.Version == "" {
		errs = append(errs, fmt.Errorf("Package version is empty"))
	} else {
		_, err := semver.ParseRange(pd.Version)
		if err != nil {
			errs = append(errs, fmt.Errorf("Invalid semver format version"))
		}
	}
	return errs
}

// GetDependencies returns the list of dependency
func (d *DependenciesFile) GetDependencies() []*Dependency {
	var dependencies []*Dependency
	for _, item := range d.Dependencies {
		dep := item
		dependencies = append(dependencies, &dep)
	}
	return dependencies
}

// GetType returns the type of dependency
func (e *Dependency) GetType() string {
	return e.Type
}

// GetTypeValue returns the dependency object that is converted
// from value string
func (e *Dependency) GetTypeValue() interface{} {
	switch e.GetType() {
	case GVKType:
		dep := GVKDependency{}
		err := json.Unmarshal([]byte(e.GetValue()), &dep)
		if err != nil {
			return nil
		}
		return dep
	case PackageType:
		dep := PackageDependency{}
		err := json.Unmarshal([]byte(e.GetValue()), &dep)
		if err != nil {
			return nil
		}
		return dep
	case LabelType:
		dep := LabelDependency{}
		err := json.Unmarshal([]byte(e.GetValue()), &dep)
		if err != nil {
			return nil
		}
		return dep
	}
	return nil
}

// GetValue returns the value content of dependency
func (e *Dependency) GetValue() string {
	return string(e.Value)
}

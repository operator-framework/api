package types

import "encoding/json"

// Property defines a single piece of the public interface for a bundle. Dependencies are specified over properties.
// The Type of the property determines how to interpret the Value, but the value is treated opaquely for
// for non-first-party types.
type Property struct {
	// The type of property. This field is required.
	Type string `json:"type" yaml:"type"`

	// The serialized value of the propertuy
	Value json.RawMessage `json:"value" yaml:"value"`
}

type GVKProperty struct {
	// The group of GVK based property
	Group string `json:"group" yaml:"group"`

	// The kind of GVK based property
	Kind string `json:"kind" yaml:"kind"`

	// The version of the API
	Version string `json:"version" yaml:"version"`
}

type PackageProperty struct {
	// The name of package such as 'etcd'
	PackageName string `json:"packageName" yaml:"packageName"`

	// The version of package in semver format
	Version string `json:"version" yaml:"version"`
}

type DeprecatedProperty struct {
	// Whether the bundle is deprecated
}

type LabelProperty struct {
	// The name of Label
	Label string `json:"label" yaml:"label"`
}
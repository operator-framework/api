package metadata

import (
	"reflect"

	"sigs.k8s.io/yaml"
)

// AnnotationsFile holds annotation information about a bundle
type AnnotationsFile struct {
	// Annotations is a set of annotations for a given bundle with both tagged and arbitrary fields.
	// Tagged fields are versioned based on the type of Annotations.
	Annotations AnnotationsV1 `json:"annotations" yaml:"annotations"`
}

// annotationsFileInternal is used for (un)marshaling an AnnotationsFile, which
// contains both tagged and arbitrary fields.
type annotationsFileInternal struct {
	Annotations map[string]string `json:"annotations"`
}

func (af AnnotationsFile) Marshal() ([]byte, error) {
	b, err := yaml.Marshal(af)
	if err != nil {
		return nil, err
	}

	// Add arbitrary annotations to the set of all annotations in af.
	annotations := make(map[string]string)
	if err := yaml.Unmarshal(b, annotations); err != nil {
		return nil, err
	}
	for label, value := range af.Annotations.AnnotationsRaw {
		annotations[label] = value
	}

	return yaml.Marshal(annotationsFileInternal{annotations})
}

func (af *AnnotationsFile) Unmarshal(b []byte) error {
	if err := yaml.Unmarshal(b, af); err != nil {
		return err
	}

	// Add arbitrary annotations in b to the set of raw annotations in af,
	// since the YAML unmarshaler cannot inherently do this.
	raw := annotationsFileInternal{make(map[string]string)}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	for _, label := range af.Annotations.GetLabels() {
		delete(raw.Annotations, label)
	}
	af.Annotations.AnnotationsRaw = make(map[string]string)
	for label, value := range raw.Annotations {
		af.Annotations.AnnotationsRaw[label] = value
	}

	return nil
}

// IsEmpty returns true if af is empty.
func (af AnnotationsFile) IsEmpty() bool {
	// A non-empty AnnotationsRaw value means af is not the zero value.
	if len(af.Annotations.AnnotationsRaw) != 0 {
		return false
	}
	// Value.IsZero() returns true if the map was initialized, which isn't correct for our purposes.
	af.Annotations.AnnotationsRaw = nil
	return reflect.ValueOf(af).IsZero()
}

// AnnotationsV1 is a list of annotations for a given bundle.
type AnnotationsV1 struct {
	// MediaType is this bundle's media type. Valid values are:
	// - "registry+v1": operator-registry manifests.
	// - "helm": Helm charts.
	// - "plain": standard Kubernetes manifests.
	MediaType string `json:"operators.operatorframework.io.bundle.mediatype.v1" yaml:"operators.operatorframework.io.bundle.mediatype.v1"`

	// ManifestsDir contains the relative manifests directory path. This path is relative
	// to the parent directory of the metadata dir, i.e. filepath.Dir(MetadataDir).
	ManifestsDir string `json:"operators.operatorframework.io.bundle.manifests.v1" yaml:"operators.operatorframework.io.bundle.manifests.v1"`

	// MetadataDir contains the relative metadata directory path. This path is relative
	// to the bundle root directory, i.e. filepath.Dir(<bundle root>).
	MetadataDir string `json:"operators.operatorframework.io.bundle.metadata.v1" yaml:"operators.operatorframework.io.bundle.metadata.v1"`

	// PackageName is the name of the overall package, ala `etcd`.
	PackageName string `json:"operators.operatorframework.io.bundle.package.v1" yaml:"operators.operatorframework.io.bundle.package.v1"`

	// Channels are a comma separated list of the declared channels for the bundle, ala `stable` or `alpha`.
	Channels string `json:"operators.operatorframework.io.bundle.channels.v1" yaml:"operators.operatorframework.io.bundle.channels.v1"`

	// DefaultChannelName is, if specified, the name of the default channel for the package. The
	// default channel will be installed if no other channel is explicitly given. If the package
	// has a single channel, then that channel is implicitly the default.
	DefaultChannelName string `json:"operators.operatorframework.io.bundle.channel.default.v1" yaml:"operators.operatorframework.io.bundle.channel.default.v1"`

	// AnnotationsRaw contains all other annotations in an annotations file
	// that do not have keys that match the above tags.
	AnnotationsRaw map[string]string `json:"-" yaml:"-"`
}

func (a AnnotationsV1) GetLabels() []string {
	return []string{
		"operators.operatorframework.io.bundle.mediatype.v1",
		"operators.operatorframework.io.bundle.manifests.v1",
		"operators.operatorframework.io.bundle.metadata.v1",
		"operators.operatorframework.io.bundle.package.v1",
		"operators.operatorframework.io.bundle.channels.v1",
		"operators.operatorframework.io.bundle.channel.default.v1",
	}
}

// DependenciesFile holds dependency information about a bundle
type DependenciesFile struct {
	// Dependencies is a list of dependencies for a given bundle
	Dependencies []Dependency `json:"dependencies" yaml:"dependencies"`
}

// Dependencies is a list of dependencies for a given bundle
type Dependency struct {
	// The type of dependency. It can be `olm.package` for operator-version based
	// dependency or `olm.gvk` for gvk based dependency. This field is required.
	Type string `json:"type" yaml:"type"`

	// The value of the dependency (either GVKDependency or PackageDependency)
	Value string `json:"value" yaml:"value"`
}

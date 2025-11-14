package release

import (
	"encoding/json"
	"strings"

	semver "github.com/blang/semver/v4"
)

// +k8s:openapi-gen=true
// OperatorRelease is a wrapper around a slice of semver.PRVersion which supports correct
// marshaling to YAML and JSON.
// +kubebuilder:validation:Type=string
type OperatorRelease struct {
	Release []semver.PRVersion `json:"-"`
}

// DeepCopyInto creates a deep-copy of the Version value.
func (v *OperatorRelease) DeepCopyInto(out *OperatorRelease) {
	out.Release = make([]semver.PRVersion, len(v.Release))
	copy(out.Release, v.Release)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (v OperatorRelease) MarshalJSON() ([]byte, error) {
	segments := []string{}
	for _, segment := range v.Release {
		segments = append(segments, segment.String())
	}
	return json.Marshal(strings.Join(segments, "."))
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (v *OperatorRelease) UnmarshalJSON(data []byte) (err error) {
	var versionString string

	if err = json.Unmarshal(data, &versionString); err != nil {
		return
	}

	segments := strings.Split(versionString, ".")
	for _, segment := range segments {
		release, err := semver.NewPRVersion(segment)
		if err != nil {
			return err
		}
		v.Release = append(v.Release, release)
	}

	return nil
}

// OpenAPISchemaType is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
//
// See: https://github.com/kubernetes/kube-openapi/tree/master/pkg/generators
func (_ OperatorRelease) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
// "semver" is not a standard openapi format but tooling may use the value regardless
func (_ OperatorRelease) OpenAPISchemaFormat() string { return "semver" }

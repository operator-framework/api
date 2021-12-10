package constraints

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type spec struct {
		name          string
		input         json.RawMessage
		expConstraint Constraint
		expError      string
	}

	specs := []spec{
		{
			name:  "Valid/BasicGVK",
			input: json.RawMessage(inputBasicGVK),
			expConstraint: Constraint{
				Message: "blah",
				GVK:     &GVKConstraint{Group: "example.com", Version: "v1", Kind: "Foo"},
			},
		},
		{
			name:  "Valid/BasicPackage",
			input: json.RawMessage(inputBasicPackage),
			expConstraint: Constraint{
				Message: "blah",
				Package: &PackageConstraint{PackageName: "foo", VersionRange: ">=1.0.0"},
			},
		},
		{
			name:  "Valid/BasicAll",
			input: json.RawMessage(fmt.Sprintf(inputBasicCompoundTmpl, "all")),
			expConstraint: Constraint{
				Message: "blah",
				All: &CompoundConstraint{
					Constraints: []Constraint{
						{
							Message: "blah blah",
							Package: &PackageConstraint{PackageName: "fuz", VersionRange: ">=1.0.0"},
						},
					},
				},
			},
		},
		{
			name:  "Valid/BasicAny",
			input: json.RawMessage(fmt.Sprintf(inputBasicCompoundTmpl, "any")),
			expConstraint: Constraint{
				Message: "blah",
				Any: &CompoundConstraint{
					Constraints: []Constraint{
						{
							Message: "blah blah",
							Package: &PackageConstraint{PackageName: "fuz", VersionRange: ">=1.0.0"},
						},
					},
				},
			},
		},
		{
			name:  "Valid/BasicNone",
			input: json.RawMessage(fmt.Sprintf(inputBasicCompoundTmpl, "none")),
			expConstraint: Constraint{
				Message: "blah",
				None: &CompoundConstraint{
					Constraints: []Constraint{
						{
							Message: "blah blah",
							Package: &PackageConstraint{PackageName: "fuz", VersionRange: ">=1.0.0"},
						},
					},
				},
			},
		},
		{
			name:  "Valid/Complex",
			input: json.RawMessage(inputComplex),
			expConstraint: Constraint{
				Message: "blah",
				All: &CompoundConstraint{
					Constraints: []Constraint{
						{Package: &PackageConstraint{PackageName: "fuz", VersionRange: ">=1.0.0"}},
						{GVK: &GVKConstraint{Group: "fals.example.com", Kind: "Fal", Version: "v1"}},
						{
							Message: "foo and buf must be stable versions",
							All: &CompoundConstraint{
								Constraints: []Constraint{
									{Package: &PackageConstraint{PackageName: "foo", VersionRange: ">=1.0.0"}},
									{Package: &PackageConstraint{PackageName: "buf", VersionRange: ">=1.0.0"}},
									{GVK: &GVKConstraint{Group: "foos.example.com", Kind: "Foo", Version: "v1"}},
								},
							},
						},
						{
							Message: "blah blah",
							Any: &CompoundConstraint{
								Constraints: []Constraint{
									{GVK: &GVKConstraint{Group: "foos.example.com", Kind: "Foo", Version: "v1beta1"}},
									{GVK: &GVKConstraint{Group: "foos.example.com", Kind: "Foo", Version: "v1beta2"}},
									{GVK: &GVKConstraint{Group: "foos.example.com", Kind: "Foo", Version: "v1"}},
								},
							},
						},
						{
							None: &CompoundConstraint{
								Constraints: []Constraint{
									{GVK: &GVKConstraint{Group: "bazs.example.com", Kind: "Baz", Version: "v1alpha1"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Invalid/TooLarge",
			input: func(t *testing.T) json.RawMessage {
				p := make([]byte, maxConstraintSize+1)
				_, err := rand.Read(p)
				require.NoError(t, err)
				return json.RawMessage(p)
			}(t),
			expError: ErrMaxConstraintSizeExceeded.Error(),
		},
		{
			name: "Invalid/UnknownField",
			input: json.RawMessage(
				`{"message": "something", "arbitrary": {"key": "value"}}`,
			),
			expError: `json: unknown field "arbitrary"`,
		},
	}

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			constraint, err := Parse(s.input)
			if s.expError == "" {
				require.NoError(t, err)
				require.Equal(t, s.expConstraint, constraint)
			} else {
				require.EqualError(t, err, s.expError)
			}
		})
	}
}

const (
	inputBasicGVK = `{
		"message": "blah",
		"gvk": {
			"group": "example.com",
			"version": "v1",
			"kind": "Foo"
		}
	}`

	inputBasicPackage = `{
		"message": "blah",
		"package": {
			"packageName": "foo",
			"versionRange": ">=1.0.0"
		}
	}`

	inputBasicCompoundTmpl = `{
"message": "blah",
"%s": {
	"constraints": [
		{
			"message": "blah blah",
			"package": {
				"packageName": "fuz",
				"versionRange": ">=1.0.0"
			}
		}
	]
}}
`

	inputComplex = `{
"message": "blah",
"all": {
	"constraints": [
		{
			"package": {
				"packageName": "fuz",
				"versionRange": ">=1.0.0"
			}
		},
		{
			"gvk": {
				"group": "fals.example.com",
				"version": "v1",
				"kind": "Fal"
			}
		},
		{
			"message": "foo and buf must be stable versions",
			"all": {
				"constraints": [
					{
						"package": {
							"packageName": "foo",
							"versionRange": ">=1.0.0"
						}
					},
					{
						"package": {
							"packageName": "buf",
							"versionRange": ">=1.0.0"
						}
					},
					{
						"gvk": {
							"group": "foos.example.com",
							"version": "v1",
							"kind": "Foo"
						}
					}
				]
			}
		},
		{
			"message": "blah blah",
			"any": {
				"constraints": [
					{
						"gvk": {
							"group": "foos.example.com",
							"version": "v1beta1",
							"kind": "Foo"
						}
					},
					{
						"gvk": {
							"group": "foos.example.com",
							"version": "v1beta2",
							"kind": "Foo"
						}
					},
					{
						"gvk": {
							"group": "foos.example.com",
							"version": "v1",
							"kind": "Foo"
						}
					}
				]
			}
		},
		{
			"none": {
				"constraints": [
					{
						"gvk": {
							"group": "bazs.example.com",
							"version": "v1alpha1",
							"kind": "Baz"
						}
					}
				]
			}
		}
	]
}}
`
)

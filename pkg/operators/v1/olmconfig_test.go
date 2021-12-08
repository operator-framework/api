package v1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func boolPointer(in bool) *bool {
	return &in
}

func TestCopiedCSVsAreEnabled(t *testing.T) {
	tests := []struct {
		description string
		olmConfig   *OLMConfig
		expected    bool
	}{
		{
			description: "NilConfig",
			olmConfig:   nil,
			expected:    true,
		},
		{
			description: "MissingSpec",
			olmConfig:   &OLMConfig{},
			expected:    true,
		},
		{
			description: "MissingFeatures",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{},
			},
			expected: true,
		},
		{
			description: "MissingDisableCopiedCSVs",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{},
			},
			expected: true,
		},
		{
			description: "CopiedCSVsDisabled",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{
					Features: &Features{
						DisableCopiedCSVs: boolPointer(true),
					},
				},
			},
			expected: false,
		},
		{
			description: "CopiedCSVsEnabled",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{
					Features: &Features{
						DisableCopiedCSVs: boolPointer(false),
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.olmConfig.CopiedCSVsAreEnabled())
		})
	}
}

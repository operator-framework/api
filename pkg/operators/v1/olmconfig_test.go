package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func boolPointer(in bool) *bool {
	return &in
}
func TestPackageServerSyncInterval(t *testing.T) {
	five := time.Minute * 5
	one := time.Second * 60

	fiveParsed, err := time.ParseDuration("5m")
	require.NoError(t, err)

	oneParsed, err := time.ParseDuration("60s")
	require.NoError(t, err)

	tests := []struct {
		description string
		olmConfig   *OLMConfig
		expected    *time.Duration
	}{
		{
			description: "NilConfig",
			olmConfig:   nil,
			expected:    nil,
		},
		{
			description: "MissingSpec",
			olmConfig:   &OLMConfig{},
			expected:    nil,
		},
		{
			description: "MissingFeatures",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{},
			},
			expected: nil,
		},
		{
			description: "MissingPackageServerInterval",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{},
			},
			expected: nil,
		},
		{
			description: "PackageServerInterval5m",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{
					Features: &Features{
						PackageServerSyncInterval: &metav1.Duration{Duration: fiveParsed},
					},
				},
			},
			expected: &five,
		},
		{
			description: "PackageServerInterval60s",
			olmConfig: &OLMConfig{
				Spec: OLMConfigSpec{
					Features: &Features{
						PackageServerSyncInterval: &metav1.Duration{Duration: oneParsed},
					},
				},
			},
			expected: &one,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.olmConfig.PackageServerSyncInterval())
		})
	}
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

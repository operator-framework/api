package v1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgradeStrategy(t *testing.T) {
	tests := []struct {
		description string
		og          *OperatorGroup
		expected    UpgradeStrategy
	}{
		{
			description: "NoSpec",
			og:          &OperatorGroup{},
			expected:    UpgradeStrategyDefault,
		},
		{
			description: "NoUpgradeStrategy",
			og: &OperatorGroup{
				Spec: OperatorGroupSpec{},
			},
			expected: UpgradeStrategyDefault,
		},
		{
			description: "NoUpgradeStrategy",
			og: &OperatorGroup{
				Spec: OperatorGroupSpec{
					UpgradeStrategy: "",
				},
			},
			expected: UpgradeStrategyDefault,
		},
		{
			description: "NonSupportedUpgradeStrategy",
			og: &OperatorGroup{
				Spec: OperatorGroupSpec{
					UpgradeStrategy: "foo",
				},
			},
			expected: UpgradeStrategyDefault,
		},
		{
			description: "DefaultUpgradeStrategy",
			og: &OperatorGroup{
				Spec: OperatorGroupSpec{
					UpgradeStrategy: "Default",
				},
			},
			expected: UpgradeStrategyDefault,
		},
		{
			description: "UnsafeFailForwardUpgradeStrategy",
			og: &OperatorGroup{
				Spec: OperatorGroupSpec{
					UpgradeStrategy: "TechPreviewUnsafeFailForward",
				},
			},
			expected: UpgradeStrategyUnsafeFailForward,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.og.UpgradeStrategy())
		})
	}
}

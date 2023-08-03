package v1alpha1

import (
	"fmt"
	"testing"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestSetRequirementStatus(t *testing.T) {
	csv := ClusterServiceVersion{}
	status := []RequirementStatus{{Group: "test", Version: "test", Kind: "Test", Name: "test", Status: "test", UUID: "test"}}
	csv.SetRequirementStatus(status)
	require.Equal(t, csv.Status.RequirementStatus, status)
}

func TestSetPhase(t *testing.T) {
	tests := []struct {
		currentPhase      ClusterServiceVersionPhase
		currentConditions []ClusterServiceVersionCondition
		inPhase           ClusterServiceVersionPhase
		outPhase          ClusterServiceVersionPhase
		description       string
	}{
		{
			currentPhase:      "",
			currentConditions: []ClusterServiceVersionCondition{},
			inPhase:           CSVPhasePending,
			outPhase:          CSVPhasePending,
			description:       "NoPhase",
		},
		{
			currentPhase:      CSVPhasePending,
			currentConditions: []ClusterServiceVersionCondition{{Phase: CSVPhasePending}},
			inPhase:           CSVPhasePending,
			outPhase:          CSVPhasePending,
			description:       "SamePhase",
		},
		{
			currentPhase:      CSVPhasePending,
			currentConditions: []ClusterServiceVersionCondition{{Phase: CSVPhasePending}},
			inPhase:           CSVPhaseInstalling,
			outPhase:          CSVPhaseInstalling,
			description:       "DifferentPhase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			csv := ClusterServiceVersion{
				Status: ClusterServiceVersionStatus{
					Phase:      tt.currentPhase,
					Conditions: tt.currentConditions,
				},
			}
			now := metav1.Now()
			csv.SetPhase(tt.inPhase, "test", "test", &now)
			require.EqualValues(t, tt.outPhase, csv.Status.Phase)
		})
	}
}

func TestIsObsolete(t *testing.T) {
	tests := []struct {
		currentPhase      ClusterServiceVersionPhase
		currentConditions []ClusterServiceVersionCondition
		out               bool
		description       string
	}{
		{
			currentPhase:      "",
			currentConditions: []ClusterServiceVersionCondition{},
			out:               false,
			description:       "NoPhase",
		},
		{
			currentPhase:      CSVPhasePending,
			currentConditions: []ClusterServiceVersionCondition{{Phase: CSVPhasePending}},
			out:               false,
			description:       "Pending",
		},
		{
			currentPhase:      CSVPhaseReplacing,
			currentConditions: []ClusterServiceVersionCondition{{Phase: CSVPhaseReplacing, Reason: CSVReasonBeingReplaced}},
			out:               true,
			description:       "Replacing",
		},
		{
			currentPhase:      CSVPhaseDeleting,
			currentConditions: []ClusterServiceVersionCondition{{Phase: CSVPhaseDeleting, Reason: CSVReasonReplaced}},
			out:               true,
			description:       "CSVPhaseDeleting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			csv := ClusterServiceVersion{
				Status: ClusterServiceVersionStatus{
					Phase:      tt.currentPhase,
					Conditions: tt.currentConditions,
				},
			}
			require.Equal(t, csv.IsObsolete(), tt.out)
		})
	}
}

func TestSupports(t *testing.T) {
	tests := []struct {
		description       string
		installModeSet    InstallModeSet
		operatorNamespace string
		namespaces        []string
		expectedErr       error
	}{
		{
			description: "NoNamespaces",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{},
			expectedErr:       fmt.Errorf("operatorgroup has invalid selected namespaces, cannot configure to watch zero namespaces"),
		},
		{
			description: "OwnNamespace/OperatorNamespace/Supported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: false,
				InstallModeTypeMultiNamespace:  false,
				InstallModeTypeAllNamespaces:   false,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"operators"},
			expectedErr:       nil,
		},
		{
			description: "SingleNamespace/OtherNamespace/Supported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    false,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  false,
				InstallModeTypeAllNamespaces:   false,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0"},
			expectedErr:       nil,
		},
		{
			description: "MultiNamespace/OtherNamespaces/Supported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    false,
				InstallModeTypeSingleNamespace: false,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   false,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0", "ns-2"},
			expectedErr:       nil,
		},
		{
			description: "AllNamespaces/NamespaceAll/Supported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    false,
				InstallModeTypeSingleNamespace: false,
				InstallModeTypeMultiNamespace:  false,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{""},
			expectedErr:       nil,
		},
		{
			description: "OwnNamespace/OperatorNamespace/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    false,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"operators"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch own namespace", InstallModeTypeOwnNamespace),
		},
		{
			description: "OwnNamespace/IncludesOperatorNamespace/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    false,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0", "operators"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch own namespace", InstallModeTypeOwnNamespace),
		},
		{
			description: "MultiNamespace/OtherNamespaces/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  false,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0", "ns-1"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch 2 namespaces", InstallModeTypeMultiNamespace),
		},
		{
			description: "SingleNamespace/OtherNamespace/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: false,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch one namespace", InstallModeTypeSingleNamespace),
		},
		{
			description: "AllNamespaces/NamespaceAll/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   false,
			},
			operatorNamespace: "operators",
			namespaces:        []string{""},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch all namespaces", InstallModeTypeAllNamespaces),
		},
		{
			description: "AllNamespaces/IncludingNamespaceAll/Unsupported",
			installModeSet: InstallModeSet{
				InstallModeTypeOwnNamespace:    true,
				InstallModeTypeSingleNamespace: true,
				InstallModeTypeMultiNamespace:  true,
				InstallModeTypeAllNamespaces:   true,
			},
			operatorNamespace: "operators",
			namespaces:        []string{"", "ns-0"},
			expectedErr:       fmt.Errorf("operatorgroup has invalid selected namespaces, NamespaceAll found when |selected namespaces| > 1"),
		},
		{
			description:       "NoNamespaces/EmptyInstallModeSet/Unsupported",
			installModeSet:    InstallModeSet{},
			operatorNamespace: "",
			namespaces:        []string{},
			expectedErr:       fmt.Errorf("operatorgroup has invalid selected namespaces, cannot configure to watch zero namespaces"),
		},
		{
			description:       "MultiNamespace/OtherNamespaces/EmptyInstallModeSet/Unsupported",
			installModeSet:    InstallModeSet{},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0", "ns-1"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch 2 namespaces", InstallModeTypeMultiNamespace),
		},
		{
			description:       "SingleNamespace/OtherNamespace/EmptyInstallModeSet/Unsupported",
			installModeSet:    InstallModeSet{},
			operatorNamespace: "operators",
			namespaces:        []string{"ns-0"},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch one namespace", InstallModeTypeSingleNamespace),
		},
		{
			description:       "AllNamespaces/NamespaceAll/EmptyInstallModeSet/Unsupported",
			installModeSet:    InstallModeSet{},
			operatorNamespace: "operators",
			namespaces:        []string{corev1.NamespaceAll},
			expectedErr:       fmt.Errorf("%s InstallModeType not supported, cannot configure to watch all namespaces", InstallModeTypeAllNamespaces),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.installModeSet.Supports(tt.operatorNamespace, tt.namespaces)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestSetPhaseWithConditions(t *testing.T) {
	tests := []struct {
		description   string
		limit         int
		currentLength int
		startIndex    int
	}{
		{
			// The original list is already at limit (length == limit).
			// We expect the oldest element ( item at 0 index) to be removed.
			description:   "TestSetPhaseWithConditionsLengthAtLimit",
			limit:         ConditionsLengthLimit,
			currentLength: ConditionsLengthLimit,

			// The first element from the original list should be dropped from
			// the new list.
			startIndex: 1,
		},
		{
			// The original list is 1 length away from limit.
			// We don't expect the list to be trimmed.
			description:   "TestSetPhaseWithConditionsLengthBelowLimit",
			limit:         ConditionsLengthLimit,
			currentLength: ConditionsLengthLimit - 1,

			// Everything in the original list should be preserved.
			startIndex: 0,
		},
		{
			// The original list has N more element(s) than allowed limit.
			// We expect (N + 1) oldest elements to be deleted to keep the list
			// at limit.
			description:   "TestSetPhaseWithConditionsLimitExceeded",
			limit:         ConditionsLengthLimit,
			currentLength: ConditionsLengthLimit + 10,

			// The first 11 (N=10 plus 1 to make room for the newly added
			// condition) elements from the original list should be dropped.
			startIndex: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			csv := ClusterServiceVersion{}
			csv.Status.Conditions = helperNewConditions(tt.currentLength)

			now := metav1.Now()

			oldConditionsWant := csv.Status.Conditions[tt.startIndex:]
			lastAddedConditionWant := ClusterServiceVersionCondition{
				Phase:              ClusterServiceVersionPhase("Pending"),
				LastTransitionTime: &now,
				LastUpdateTime:     &now,
				Message:            "message",
				Reason:             ConditionReason("reason"),
			}

			csv.SetPhase("Pending", "reason", "message", &now)

			conditionsGot := csv.Status.Conditions
			assert.Equal(t, tt.limit, len(conditionsGot))

			oldConditionsGot := conditionsGot[0 : len(conditionsGot)-1]
			assert.EqualValues(t, oldConditionsWant, oldConditionsGot)

			lastAddedConditionGot := conditionsGot[len(conditionsGot)-1]
			assert.Equal(t, lastAddedConditionWant, lastAddedConditionGot)
		})
	}
}

func TestWebhookDescGetValidatingConfigurations(t *testing.T) {
	expectedPort := int32(444)
	timeout := int32(32)
	webhookPath := "/test"
	failurePolicy := admissionregistrationv1.Fail
	matchPolicy := admissionregistrationv1.Exact
	sideEffect := admissionregistrationv1.SideEffectClassNone
	webhookDesc := WebhookDescription{
		GenerateName:            "foo-webhook",
		Type:                    ValidatingAdmissionWebhook,
		DeploymentName:          "foo-deployment",
		ContainerPort:           expectedPort,
		AdmissionReviewVersions: []string{"v1beta1", "v1"},
		SideEffects:             &sideEffect,
		MatchPolicy:             &matchPolicy,
		FailurePolicy:           &failurePolicy,
		ObjectSelector:          &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
		TimeoutSeconds:          &timeout,
		WebhookPath:             &webhookPath,
		Rules: []admissionregistrationv1.RuleWithOperations{
			admissionregistrationv1.RuleWithOperations{
				Operations: []admissionregistrationv1.OperationType{},
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{"*"},
					APIVersions: []string{"*"},
					Resources:   []string{"*"},
				},
			},
		},
	}
	vWebhookConfig := webhookDesc.GetValidatingWebhook("foo", nil, nil)
	require.Equal(t, expectedPort, *vWebhookConfig.ClientConfig.Service.Port)
	require.Equal(t, webhookDesc.Rules, vWebhookConfig.Rules)
	require.Equal(t, webhookDesc.FailurePolicy, vWebhookConfig.FailurePolicy)
	require.Equal(t, webhookDesc.MatchPolicy, vWebhookConfig.MatchPolicy)
	require.Equal(t, webhookDesc.ObjectSelector, vWebhookConfig.ObjectSelector)
	require.Equal(t, webhookDesc.SideEffects, vWebhookConfig.SideEffects)
	require.Equal(t, webhookDesc.TimeoutSeconds, vWebhookConfig.TimeoutSeconds)
	require.Equal(t, webhookDesc.AdmissionReviewVersions, vWebhookConfig.AdmissionReviewVersions)
	require.Equal(t, webhookDesc.WebhookPath, vWebhookConfig.ClientConfig.Service.Path)

	mWebhookConfig := webhookDesc.GetMutatingWebhook("foo", nil, nil)
	require.Equal(t, expectedPort, *mWebhookConfig.ClientConfig.Service.Port)
	require.Equal(t, webhookDesc.Rules, mWebhookConfig.Rules)
	require.Equal(t, webhookDesc.FailurePolicy, mWebhookConfig.FailurePolicy)
	require.Equal(t, webhookDesc.MatchPolicy, mWebhookConfig.MatchPolicy)
	require.Equal(t, webhookDesc.ObjectSelector, mWebhookConfig.ObjectSelector)
	require.Equal(t, webhookDesc.SideEffects, mWebhookConfig.SideEffects)
	require.Equal(t, webhookDesc.TimeoutSeconds, mWebhookConfig.TimeoutSeconds)
	require.Equal(t, webhookDesc.AdmissionReviewVersions, mWebhookConfig.AdmissionReviewVersions)
	require.Equal(t, webhookDesc.ReinvocationPolicy, mWebhookConfig.ReinvocationPolicy)
	require.Equal(t, webhookDesc.WebhookPath, mWebhookConfig.ClientConfig.Service.Path)
}

func helperNewConditions(count int) []ClusterServiceVersionCondition {
	conditions := make([]ClusterServiceVersionCondition, 0)

	for i := 1; i <= count; i++ {
		now := metav1.Now()
		condition := ClusterServiceVersionCondition{
			Phase:              ClusterServiceVersionPhase(fmt.Sprintf("phase-%d", i)),
			LastTransitionTime: &now,
			LastUpdateTime:     &now,
			Message:            fmt.Sprintf("message-%d", i),
			Reason:             ConditionReason(fmt.Sprintf("reason-%d", i)),
		}
		conditions = append(conditions, condition)
	}

	return conditions
}

func TestIsCopied(t *testing.T) {
	var testCases = []struct {
		name     string
		input    metav1.Object
		expected bool
	}{
		{
			name:     "no labels or annotations",
			input:    &metav1.ObjectMeta{},
			expected: false,
		},
		{
			name: "no labels, has annotations but missing operatorgroup namespace annotation",
			input: &metav1.ObjectMeta{
				Annotations: map[string]string{},
			},
			expected: false,
		},
		{
			name: "no labels, has operatorgroup namespace annotation matching self",
			input: &metav1.ObjectMeta{
				Namespace: "whatever",
				Annotations: map[string]string{
					"olm.operatorNamespace": "whatever",
				},
			},
			expected: false,
		},
		{
			name: "no labels, has operatorgroup namespace annotation not matching self",
			input: &metav1.ObjectMeta{
				Namespace: "whatever",
				Annotations: map[string]string{
					"olm.operatorNamespace": "other",
				},
			},
			expected: true,
		},
		{
			name: "no annotations, labels missing copied key",
			input: &metav1.ObjectMeta{
				Labels: map[string]string{},
			},
			expected: false,
		},
		{
			name: "no annotations, labels has copied key",
			input: &metav1.ObjectMeta{
				Labels: map[string]string{
					"olm.copiedFrom": "whatever",
				},
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got, expected := IsCopied(testCase.input), testCase.expected; got != expected {
				t.Errorf("got %v, expected %v", got, expected)
			}
		})
	}
}

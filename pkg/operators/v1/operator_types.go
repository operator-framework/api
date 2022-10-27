package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperatorSpec defines the desired state of Operator
type OperatorSpec struct{}

// OperatorStatus defines the observed state of an Operator and its components
type OperatorStatus struct {
	// Components describes resources that compose the operator.
	// +optional
	Components *Components `json:"components,omitempty"`
}

// ConditionType codifies a condition's type.
type ConditionType string

// Condition represent the latest available observations of an component's state.
type Condition struct {
	// Type of condition.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Last time the condition was probed
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
}

// Components tracks the resources that compose an operator.
type Components struct {
	// LabelSelector is a label query over a set of resources used to select the operator's components
	LabelSelector *metav1.LabelSelector `json:"labelSelector"`
	// Refs are a set of references to the operator's component resources, selected with LabelSelector.
	// +optional
	Refs []RichReference `json:"refs,omitempty"`
}

// Len returns the number of Refs in the component.Refs array
// Used to implement the interface required by sort.Sort function.
func (c Components) Len() int {
	return len(c.Refs)
}

// Less returns true if argument i should appear in an ordered list
// of references before argument j.
// Used to implement the sort.Sort interface.
func (c Components) Less(i, j int) bool {
	if c.Refs[i].Kind != c.Refs[j].Kind {
		return c.Refs[i].Kind < c.Refs[j].Kind
	}

	if c.Refs[i].APIVersion != c.Refs[j].APIVersion {
		return c.Refs[i].APIVersion < c.Refs[j].APIVersion
	}

	if c.Refs[i].Namespace != c.Refs[j].Namespace {
		return c.Refs[i].Namespace < c.Refs[j].Namespace
	}
	return c.Refs[i].Name < c.Refs[j].Name
}

// Swap swaps the elements with indexes i and j.
// Used to implement the sort.Sort interface.
func (c Components) Swap(i, j int) {
	c.Refs[i], c.Refs[j] = c.Refs[j], c.Refs[i]
}

// RichReference is a reference to a resource, enriched with its status conditions.
type RichReference struct {
	*corev1.ObjectReference `json:",inline"`
	// Conditions represents the latest state of the component.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=olm,scope=Cluster
// +kubebuilder:subresource:status

// Operator represents a cluster operator.
type Operator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorSpec   `json:"spec,omitempty"`
	Status OperatorStatus `json:"status,omitempty"`
}

// +genclient:nonNamespaced
// +kubebuilder:object:root=true

// OperatorList contains a list of Operators.
type OperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Operator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Operator{}, &OperatorList{})
}

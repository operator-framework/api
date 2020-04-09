package crds

import (
	"reflect"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var emptyCRD = &apiextensionsv1.CustomResourceDefinition{}

func TestGetters(t *testing.T) {
	tests := []struct {
		description string
		get         func() *apiextensionsv1.CustomResourceDefinition
	}{
		{
			description: "CatalogSource",
			get:         CatalogSource,
		},
		{
			description: "ClusterServiceVersion",
			get:         ClusterServiceVersion,
		},
		{
			description: "InstallPlan",
			get:         InstallPlan,
		},
		{
			description: "OperatorGroup",
			get:         OperatorGroup,
		},
		{
			description: "Operator",
			get:         Operator,
		},
		{
			description: "Subscription",
			get:         Subscription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			defer func() {
				if x := recover(); x != nil {
					t.Errorf("panic loading crd: %v", x)
				}
			}()

			crd := tt.get()
			if crd == nil || reflect.DeepEqual(crd, emptyCRD) {
				t.Error("loaded CustomResourceDefinition is empty")
			}
		})
	}
}

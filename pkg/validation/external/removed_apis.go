package external

import (
	"github.com/operator-framework/api/pkg/manifests"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetRemovedAPIsOn1_25From(bundle *manifests.Bundle) map[string][]string {
	deprecatedAPIs := make(map[string][]string)
	for _, obj := range bundle.Objects {
		switch usage := obj.GetObjectKind().(type) {
		case *unstructured.Unstructured:
			switch usage.GetAPIVersion() {
			case "batch/v1beta1":
				if usage.GetKind() == "CronJob" {
					addDepUsage(deprecatedAPIs, usage)
				}
			case "discovery.k8s.io/v1beta1":
				if usage.GetKind() == "EndpointSlice" {
					addDepUsage(deprecatedAPIs, usage)
				}
			case "events.k8s.io/v1beta1":
				if usage.GetKind() == "Event" {
					addDepUsage(deprecatedAPIs, usage)
				}
			case "policy/v1beta1":
				if usage.GetKind() == "PodDisruptionBudget" || usage.GetKind() == "PodSecurityPolicy" {
					addDepUsage(deprecatedAPIs, usage)
				}
			case "node.k8s.io/v1beta1":
				if usage.GetKind() == "RuntimeClass" {
					addDepUsage(deprecatedAPIs, usage)
				}
			}
		}
	}
	return deprecatedAPIs
}

func addDepUsage(deprecatedAPIs map[string][]string, u *unstructured.Unstructured) map[string][]string {
	deprecatedAPIs[u.GetKind()] = append(deprecatedAPIs[u.GetKind()], u.GetName())
	return deprecatedAPIs
}

package internal

import (
	"fmt"
	"strings"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation/errors"
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
)

// DeprecatedImages contains a list of deprecated images and their respective messages
var DeprecatedImages = map[string]string{
	"gcr.io/kubebuilder/kube-rbac-proxy": "Your bundle uses the image `gcr.io/kubebuilder/kube-rbac-proxy`. This upstream image is deprecated and may become unavailable at any point.  \n\nPlease use an equivalent image from a trusted source or update your approach to protect the metrics endpoint.  \n\nFor further information: https://github.com/kubernetes-sigs/kubebuilder/discussions/3907",
}

// ImageDeprecateValidator implements Validator to validate bundle objects for deprecated image usage.
var ImageDeprecateValidator interfaces.Validator = interfaces.ValidatorFunc(validateImageDeprecateValidator)

// validateImageDeprecateValidator checks for the presence of deprecated images in the bundle's CSV and deployment specs.
func validateImageDeprecateValidator(objs ...interface{}) (results []errors.ManifestResult) {
	for _, obj := range objs {
		switch v := obj.(type) {
		case *manifests.Bundle:
			results = append(results, validateDeprecatedImage(v))
		}
	}
	return results
}

// validateDeprecatedImage checks for deprecated images in both the CSV `RelatedImages` and deployment specs.
func validateDeprecatedImage(bundle *manifests.Bundle) errors.ManifestResult {
	result := errors.ManifestResult{}
	if bundle == nil {
		result.Add(errors.ErrInvalidBundle("Bundle is nil", nil))
		return result
	}

	result.Name = bundle.Name

	if bundle.CSV != nil {
		for _, relatedImage := range bundle.CSV.Spec.RelatedImages {
			for deprecatedImage, message := range DeprecatedImages {
				if strings.Contains(relatedImage.Image, deprecatedImage) {
					result.Add(errors.WarnFailedValidation(
						fmt.Sprintf(message, relatedImage.Image),
						relatedImage.Name,
					))
				}
			}
		}
	}

	if bundle.CSV != nil {
		for _, deploymentSpec := range bundle.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
			for _, container := range deploymentSpec.Spec.Template.Spec.Containers {
				for deprecatedImage, message := range DeprecatedImages {
					if strings.Contains(container.Image, deprecatedImage) {
						result.Add(errors.WarnFailedValidation(
							fmt.Sprintf(message, container.Image),
							deploymentSpec.Name,
						))
					}
				}
			}
		}
	}

	return result
}

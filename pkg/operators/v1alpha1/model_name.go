/*
Copyright Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

// OpenAPIModelName methods for all types that are referenced in OpenAPI schemas.
// This ensures the OpenAPI schema uses the canonical name format (com.github...)
// instead of the URL-encoded Go import path format that causes lookup failures.
//
// When Kubernetes serves OpenAPI schemas, "/" in type references gets URL-encoded
// to "~1". If the schema definition keys don't match, lookups fail with errors like:
//   unknown model in reference: "github.com~1operator-framework~1api~1..."
//
// Adding OpenAPIModelName() makes openapi-gen use the canonical format for both
// the definition key and references, ensuring they match.

// OpenAPIModelName returns the OpenAPI model name for this type.
func (APIResourceReference) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.APIResourceReference"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (APIServiceDefinitions) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.APIServiceDefinitions"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (APIServiceDescription) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.APIServiceDescription"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ActionDescriptor) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ActionDescriptor"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (AppLink) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.AppLink"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (CRDDescription) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.CRDDescription"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (CleanupSpec) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.CleanupSpec"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (CleanupStatus) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.CleanupStatus"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ClusterServiceVersionCondition) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionCondition"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ClusterServiceVersionSpec) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionSpec"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ClusterServiceVersionStatus) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionStatus"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (CustomResourceDefinitions) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.CustomResourceDefinitions"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (DependentStatus) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.DependentStatus"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (Icon) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.Icon"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (InstallMode) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.InstallMode"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (Maintainer) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.Maintainer"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (NamedInstallStrategy) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.NamedInstallStrategy"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (RelatedImage) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.RelatedImage"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (RequirementStatus) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.RequirementStatus"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ResourceInstance) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ResourceInstance"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (ResourceList) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.ResourceList"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (SpecDescriptor) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.SpecDescriptor"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (StatusDescriptor) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.StatusDescriptor"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (StrategyDeploymentPermissions) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDeploymentPermissions"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (StrategyDeploymentSpec) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDeploymentSpec"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (StrategyDetailsDeployment) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDetailsDeployment"
}

// OpenAPIModelName returns the OpenAPI model name for this type.
func (WebhookDescription) OpenAPIModelName() string {
	return "com.github.operator-framework.api.pkg.operators.v1alpha1.WebhookDescription"
}

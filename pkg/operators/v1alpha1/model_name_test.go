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

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAPIModelName(t *testing.T) {
	tests := []struct {
		name     string
		typeFunc func() string
		expected string
	}{
		{
			name:     "APIResourceReference",
			typeFunc: func() string { return APIResourceReference{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.APIResourceReference",
		},
		{
			name:     "APIServiceDefinitions",
			typeFunc: func() string { return APIServiceDefinitions{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.APIServiceDefinitions",
		},
		{
			name:     "APIServiceDescription",
			typeFunc: func() string { return APIServiceDescription{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.APIServiceDescription",
		},
		{
			name:     "ActionDescriptor",
			typeFunc: func() string { return ActionDescriptor{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ActionDescriptor",
		},
		{
			name:     "AppLink",
			typeFunc: func() string { return AppLink{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.AppLink",
		},
		{
			name:     "CRDDescription",
			typeFunc: func() string { return CRDDescription{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.CRDDescription",
		},
		{
			name:     "CleanupSpec",
			typeFunc: func() string { return CleanupSpec{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.CleanupSpec",
		},
		{
			name:     "CleanupStatus",
			typeFunc: func() string { return CleanupStatus{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.CleanupStatus",
		},
		{
			name:     "ClusterServiceVersionCondition",
			typeFunc: func() string { return ClusterServiceVersionCondition{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionCondition",
		},
		{
			name:     "ClusterServiceVersionSpec",
			typeFunc: func() string { return ClusterServiceVersionSpec{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionSpec",
		},
		{
			name:     "ClusterServiceVersionStatus",
			typeFunc: func() string { return ClusterServiceVersionStatus{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ClusterServiceVersionStatus",
		},
		{
			name:     "CustomResourceDefinitions",
			typeFunc: func() string { return CustomResourceDefinitions{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.CustomResourceDefinitions",
		},
		{
			name:     "DependentStatus",
			typeFunc: func() string { return DependentStatus{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.DependentStatus",
		},
		{
			name:     "Icon",
			typeFunc: func() string { return Icon{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.Icon",
		},
		{
			name:     "InstallMode",
			typeFunc: func() string { return InstallMode{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.InstallMode",
		},
		{
			name:     "Maintainer",
			typeFunc: func() string { return Maintainer{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.Maintainer",
		},
		{
			name:     "NamedInstallStrategy",
			typeFunc: func() string { return NamedInstallStrategy{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.NamedInstallStrategy",
		},
		{
			name:     "RelatedImage",
			typeFunc: func() string { return RelatedImage{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.RelatedImage",
		},
		{
			name:     "RequirementStatus",
			typeFunc: func() string { return RequirementStatus{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.RequirementStatus",
		},
		{
			name:     "ResourceInstance",
			typeFunc: func() string { return ResourceInstance{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ResourceInstance",
		},
		{
			name:     "ResourceList",
			typeFunc: func() string { return ResourceList{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.ResourceList",
		},
		{
			name:     "SpecDescriptor",
			typeFunc: func() string { return SpecDescriptor{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.SpecDescriptor",
		},
		{
			name:     "StatusDescriptor",
			typeFunc: func() string { return StatusDescriptor{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.StatusDescriptor",
		},
		{
			name:     "StrategyDeploymentPermissions",
			typeFunc: func() string { return StrategyDeploymentPermissions{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDeploymentPermissions",
		},
		{
			name:     "StrategyDeploymentSpec",
			typeFunc: func() string { return StrategyDeploymentSpec{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDeploymentSpec",
		},
		{
			name:     "StrategyDetailsDeployment",
			typeFunc: func() string { return StrategyDetailsDeployment{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.StrategyDetailsDeployment",
		},
		{
			name:     "WebhookDescription",
			typeFunc: func() string { return WebhookDescription{}.OpenAPIModelName() },
			expected: "com.github.operator-framework.api.pkg.operators.v1alpha1.WebhookDescription",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.typeFunc())
		})
	}
}

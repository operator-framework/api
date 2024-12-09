package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"strings"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestValidateDeprecatedImage(t *testing.T) {
	type args struct {
		bundleDir string
	}
	tests := []struct {
		name         string
		args         args
		modifyBundle func(bundle *manifests.Bundle)
		wantErr      []string
		want         map[string]string
	}{
		{
			name: "should return no error or warning for a valid bundle",
			args: args{
				bundleDir: "./testdata/valid_bundle",
			},
			modifyBundle: func(bundle *manifests.Bundle) {
				// Do not modify the bundle
			},
			wantErr: []string{},
			want:    map[string]string{},
		},
		{
			name: "should detect deprecated image in RelatedImages",
			args: args{
				bundleDir: "./testdata/valid_bundle",
			},
			modifyBundle: func(bundle *manifests.Bundle) {
				bundle.CSV.Spec.RelatedImages = append(bundle.CSV.Spec.RelatedImages, operatorsv1alpha1.RelatedImage{
					Name:  "kube-rbac-proxy",
					Image: "gcr.io/kubebuilder/kube-rbac-proxy",
				})
			},
			wantErr: []string{
				"Your bundle uses the image `gcr.io/kubebuilder/kube-rbac-proxy`",
			},
			want: map[string]string{
				"RelatedImage": "gcr.io/kubebuilder/kube-rbac-proxy",
			},
		},
		{
			name: "should detect deprecated image in DeploymentSpecs",
			args: args{
				bundleDir: "./testdata/valid_bundle",
			},
			modifyBundle: func(bundle *manifests.Bundle) {
				bundle.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs = append(
					bundle.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs,
					operatorsv1alpha1.StrategyDeploymentSpec{
						Name: "test-deployment",
						Spec: appsv1.DeploymentSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "kube-rbac-proxy",
											Image: "gcr.io/kubebuilder/kube-rbac-proxy",
										},
									},
								},
							},
						},
					},
				)
			},
			wantErr: []string{
				"Your bundle uses the image `gcr.io/kubebuilder/kube-rbac-proxy`",
			},
			want: map[string]string{
				"DeploymentSpec": "gcr.io/kubebuilder/kube-rbac-proxy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			tt.modifyBundle(bundle)
			result := validateDeprecatedImage(bundle)

			var gotWarnings []string
			for _, warn := range result.Warnings {
				gotWarnings = append(gotWarnings, warn.Error())
			}

			for _, expectedErr := range tt.wantErr {
				found := false
				for _, gotWarning := range gotWarnings {
					if strings.Contains(gotWarning, expectedErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing '%s' but did not find it", expectedErr)
				}
			}

			gotImages := make(map[string]string)
			if bundle.CSV != nil {
				for _, relatedImage := range bundle.CSV.Spec.RelatedImages {
					for deprecatedImage := range DeprecatedImages {
						if relatedImage.Image == deprecatedImage {
							gotImages["RelatedImage"] = deprecatedImage
						}
					}
				}
				for _, deploymentSpec := range bundle.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
					for _, container := range deploymentSpec.Spec.Template.Spec.Containers {
						for deprecatedImage := range DeprecatedImages {
							if container.Image == deprecatedImage {
								gotImages["DeploymentSpec"] = deprecatedImage
							}
						}
					}
				}
			}

			require.True(t, reflect.DeepEqual(tt.want, gotImages))
		})
	}
}

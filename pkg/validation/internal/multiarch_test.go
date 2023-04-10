package internal

import (
	"reflect"
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/require"
)

const opm_test_image = "quay.io/operator-framework/opm:latest"

func Test_ValidateMultiArchFrom(t *testing.T) {

	// Mock bundle
	bundleWithoutLabels, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")

	if bundleWithoutLabels.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs != nil {
		for indexDeployment, v := range bundleWithoutLabels.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
			for indexContainer, _ := range v.Spec.Template.Spec.Containers {
				bundleWithoutLabels.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs[indexDeployment].Spec.Template.Spec.Containers[indexContainer].Image = opm_test_image
			}
		}
	}

	for i, _ := range bundleWithoutLabels.CSV.Spec.RelatedImages {
		bundleWithoutLabels.CSV.Spec.RelatedImages[i].Image = opm_test_image
	}

	allLabels := map[string]string{}
	allLabels["operatorframework.io/arch.arm64"] = "supported"
	allLabels["operatorframework.io/arch.ppc64le"] = "supported"
	allLabels["operatorframework.io/arch.s390x"] = "supported"
	allLabels["operatorframework.io/arch.amd64"] = "supported"

	type args struct {
		labels    map[string]string
		bundle    *manifests.Bundle
		bundleDir string
	}
	tests := []struct {
		name        string
		args        args
		wantWarning bool
		wantError   bool
		warnStrings []string
		errStrings  []string
	}{
		{
			name: "should warning when is missing allLabels for the arch types found on the images",
			args: args{
				bundle: bundleWithoutLabels,
			},
			wantWarning: true,
			warnStrings: []string{"Warning: Value etcdoperator.v0.9.4: check if the CSV is missing the " +
				"label (operatorframework.io/arch.<value>) for the Arch(s): " +
				"[\"amd64\" \"arm64\" \"ppc64le\" \"s390x\"]. Be aware that your Operator manager " +
				"image [\"quay.io/operator-framework/opm:latest\"] provides this support. " +
				"Thus, it is very likely that you want to provide it and if you support more than " +
				"amd64 architectures, you MUST,use the required labels for all which are " +
				"supported.Otherwise, your solution cannot be listed on the cluster for these " +
				"architectures"},
		},
		{
			name: "should successfully pass when the bundle has all labels",
			args: args{
				bundle: bundleWithoutLabels,
				labels: allLabels,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.bundleDir) > 0 {
				bundleLoaded, err := manifests.GetBundleFromDir(tt.args.bundleDir)
				require.Equal(t, err, nil)
				tt.args.bundle = bundleLoaded
			}

			if len(tt.args.labels) > 0 {
				tt.args.bundle.CSV.Labels = tt.args.labels
			}

			results := validateMultiArchWith(tt.args.bundle, "")
			t.Log(results.Warnings)
			t.Log(results.Errors)

			require.Equal(t, tt.wantWarning, len(results.Warnings) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(results.Warnings))
				for _, w := range results.Warnings {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(results.Errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(results.Errors))
				for _, w := range results.Errors {
					wString := w.Error()
					require.Contains(t, tt.errStrings, wString)
				}
			}
		})
	}
}

func Test_LoadImagesFromCSV(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	type fields struct {
		Bundle        *manifests.Bundle
		InstallImages []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "should successfully get all install images from the CSV ",
			fields: fields{
				Bundle: validBundle,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &multiArchValidator{
				bundle: tt.fields.Bundle,
			}
			mb.loadImagesFromCSV()
			require.Equal(t, len(mb.managerImages), 1)
			require.Equal(t, len(mb.managerImages["quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"]), 0)
			require.Equal(t, 1, len(mb.otherCSVDeploymentImages))
			require.Equal(t, len(mb.otherCSVDeploymentImages["quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"]), 0)
			require.Equal(t, 0, len(mb.relatedImages))
		})
	}
}

func Test_LoadImagesFromCSVWithRelatedImage(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	validBundle.CSV.Spec.RelatedImages = []v1alpha1.RelatedImage{
		{Image: "related-image-test", Name: "etcd-operator"},
	}

	type fields struct {
		Bundle       *manifests.Bundle
		RelateImages []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "should successfully get all related install images from the CSV ",
			fields: fields{
				Bundle: validBundle,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &multiArchValidator{
				bundle: tt.fields.Bundle,
			}
			mb.loadImagesFromCSV()
			require.Equal(t, len(mb.managerImages), 1)
			require.Equal(t, len(mb.managerImages["quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"]), 0)
			require.Equal(t, 1, len(mb.otherCSVDeploymentImages))
			require.Equal(t, len(mb.otherCSVDeploymentImages["quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"]), 0)
			require.Equal(t, 1, len(mb.relatedImages))
			require.Equal(t, len(mb.relatedImages["related-image-test"]), 0)
		})
	}
}

func Test_ValidateContainerTool(t *testing.T) {
	type args struct {
		containerTool string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return docker when is empty",
			args: args{
				containerTool: "",
			},
			want: "docker",
		},
		{
			name: "should return docker when is none",
			args: args{
				containerTool: "",
			},
			want: "docker",
		},
		{
			name: "should return error when is not one of supported options",
			args: args{
				containerTool: "invalid",
			},
			want:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateContainerTool(tt.args.containerTool)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateContainerTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateContainerTool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RunManifestInspect(t *testing.T) {
	type args struct {
		image string
		tool  string
	}
	tests := []struct {
		name    string
		args    args
		want    manifestInspect
		wantErr bool
	}{
		{
			name: "should return the data from the manifest",
			args: args{
				tool:  "docker",
				image: opm_test_image,
			},
			want: manifestInspect{[]manifestData{
				{platform{Architecture: "amd64", OS: "linux"}},
				{platform{Architecture: "arm64", OS: "linux"}},
				{platform{Architecture: "ppc64le", OS: "linux"}},
				{platform{Architecture: "s390x", OS: "linux"}}}},
			wantErr: false,
		},
		{
			name: "should fail when is not possible to inspect",
			args: args{
				tool:  "docker",
				image: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runManifestInspect(tt.args.image, tt.args.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf("runManifestInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runManifestInspect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_LoadAllPossibleArchSupported(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	managerArchAmd64 := map[string]string{}
	managerArchAmd64["amd64"] = "amd64"

	validBundleWithInfraLabels, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	mockInfraLabelsArch := map[string]string{}
	mockInfraLabelsArch["operatorframework.io/arch.amd64"] = "supported"
	mockInfraLabelsArch["operatorframework.io/arch.ppc64le"] = "supported"
	validBundleWithInfraLabels.CSV.Labels = mockInfraLabelsArch

	managerMultiArch := map[string]string{}
	managerMultiArch["amd64"] = "amd64"
	managerMultiArch["ppc64le"] = "ppc64le"

	type fields struct {
		bundle *manifests.Bundle
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "should add amd64 when no label is informed",
			fields: fields{
				bundle: validBundle,
			},
			want: managerArchAmd64,
		},
		{
			name: "should return the infra labels when informed",
			fields: fields{
				bundle: validBundleWithInfraLabels,
			},
			want: managerMultiArch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &multiArchValidator{
				bundle: tt.fields.bundle,
			}

			data.loadImagesFromCSV()
			data.loadInfraLabelsFromCSV()
			data.loadAllPossibleArchSupported()

			if !reflect.DeepEqual(data.managerArchs, tt.want) {
				t.Errorf("loadAllPossibleArchSupported() got = %v, want %v", data.managerArchs, tt.want)
			}
		})
	}
}

func Test_LoadAllPossibleOsSupported(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	managerOsLinux := map[string]string{}
	managerOsLinux["linux"] = "linux"

	validBundleWithInfraLabels, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	mockInfraLabelsOs := map[string]string{}
	mockInfraLabelsOs["operatorframework.io/os.linux"] = "supported"
	mockInfraLabelsOs["operatorframework.io/os.other"] = "supported"
	validBundleWithInfraLabels.CSV.Labels = mockInfraLabelsOs

	managerMultiOs := map[string]string{}
	managerMultiOs["linux"] = "linux"
	managerMultiOs["other"] = "other"

	type fields struct {
		bundle *manifests.Bundle
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "should add linux when no label is informed",
			fields: fields{
				bundle: validBundle,
			},
			want: managerOsLinux,
		},
		{
			name: "should return the infra labels when informed",
			fields: fields{
				bundle: validBundleWithInfraLabels,
			},
			want: managerMultiOs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &multiArchValidator{
				bundle: tt.fields.bundle,
			}

			data.loadImagesFromCSV()
			data.loadInfraLabelsFromCSV()
			data.loadAllPossibleOsSupported()

			if !reflect.DeepEqual(data.managerOs, tt.want) {
				t.Errorf("loadAllPossibleOsSupported() got = %v, want %v", data.managerOs, tt.want)
			}
		})
	}
}

func Test_multiArchValidator_checkSupportDefined(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")

	validBundleWithLabels, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	allLabels := map[string]string{}
	allLabels["operatorframework.io/arch.arm64"] = "supported"
	allLabels["operatorframework.io/os.linux"] = "supported"
	allLabels["operatorframework.io/os.other"] = "supported"
	validBundleWithLabels.CSV.Labels = allLabels

	type fields struct {
		bundle *manifests.Bundle
	}
	tests := []struct {
		name        string
		fields      fields
		wantWarning bool
		wantError   bool
		warnStrings []string
		errStrings  []string
	}{
		{
			name: "should raise no error or warning when has the support defined",
			fields: fields{
				bundle: validBundle,
			},
		},
		{
			name: "should raise an error when all images does not have the support defined via the labels",
			fields: fields{
				bundle: validBundleWithLabels,
			},
			wantError: true,
			errStrings: []string{"not all images specified are providing the support described via the CSV labels. Note that (OS.architecture): (linux.arm64) was not found for the image(s) [quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b quay.io/coreos/etcd-operator@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b]",
				"not all images specified are providing the support described via the CSV labels. Note that (OS.architecture): (other.arm64) was not found for the image(s) [quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b quay.io/coreos/etcd-operator@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &multiArchValidator{
				bundle: tt.fields.bundle,
			}

			data.loadInfraLabelsFromCSV()
			data.loadImagesFromCSV()

			// Mock inspected platform
			for key, _ := range data.managerImages {
				data.managerImages[key] = []platform{{"linux", "amd64"}}
			}
			for key, _ := range data.otherCSVDeploymentImages {
				data.otherCSVDeploymentImages[key] = []platform{{"linux", "amd64"}}
			}

			data.loadAllPossibleArchSupported()
			data.loadAllPossibleOsSupported()

			data.checkSupportDefined()

			t.Log(data.warns)
			t.Log(data.errors)

			require.Equal(t, tt.wantWarning, len(data.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(data.warns))
				for _, w := range data.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(data.errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(data.errors))
				for _, w := range data.errors {
					wString := w.Error()
					require.Contains(t, tt.errStrings, wString)
				}
			}
		})
	}
}

func Test_multiArchValidator_checkMissingLabels(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")

	validBundleWithLabels, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	allLabels := map[string]string{}
	allLabels["operatorframework.io/arch.amd64"] = "supported"
	allLabels["operatorframework.io/arch.arm64"] = "supported"
	allLabels["operatorframework.io/os.linux"] = "supported"
	validBundleWithLabels.CSV.Labels = allLabels

	type fields struct {
		bundle             *manifests.Bundle
		supportedPlatforms []platform
	}
	tests := []struct {
		name        string
		fields      fields
		wantWarning bool
		wantError   bool
		warnStrings []string
		errStrings  []string
	}{
		{
			name: "should raise no error or warning when only supports linux.amd64 (no labels are required)",
			fields: fields{
				bundle:             validBundle,
				supportedPlatforms: []platform{{"linux", "amd64"}},
			},
		},
		{
			name: "should not raise an error when has all labels for all that is supported is set",
			fields: fields{
				bundle: validBundleWithLabels,
				supportedPlatforms: []platform{
					{"other", "amd64"},
					{"other", "arm64"},
					{"linux", "amd64"},
					{"linux", "arm64"},
				},
			},
		},
		{
			name: "should raise a warning when is missing a label",
			fields: fields{
				bundle: validBundleWithLabels,
				supportedPlatforms: []platform{
					{"other", "amd64"},
					{"other", "arm64"},
					{"other", "missing"},
					{"linux", "amd64"},
					{"linux", "arm64"},
					{"linux", "missing"},
				},
			},
			wantWarning: true,
			warnStrings: []string{"check if the CSV is missing the label (operatorframework.io/arch.<value>) for the Arch(s): [\"missing\"]. Be aware that your Operator manager image [\"quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b\"] provides this support. Thus, it is very likely that you want to provide it and if you support more than amd64 architectures, you MUST,use the required labels for all which are supported.Otherwise, your solution cannot be listed on the cluster for these architectures"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &multiArchValidator{
				bundle: tt.fields.bundle,
			}

			data.loadInfraLabelsFromCSV()
			data.loadImagesFromCSV()

			// Mock inspected platform
			for key, _ := range data.managerImages {
				data.managerImages[key] = tt.fields.supportedPlatforms
			}

			data.loadAllPossibleArchSupported()
			data.loadAllPossibleOsSupported()
			data.checkMissingLabelsForArchs()

			t.Log(data.warns)
			t.Log(data.errors)

			require.Equal(t, tt.wantWarning, len(data.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(data.warns))
				for _, w := range data.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(data.errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(data.errors))
				for _, w := range data.errors {
					wString := w.Error()
					require.Contains(t, tt.errStrings, wString)
				}
			}
		})
	}
}

func Test_multiArchValidator_checkNodeAffinity(t *testing.T) {
	validBundle, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")

	validBundleMissingNodeAffinity, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	for i := range validBundleMissingNodeAffinity.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
		validBundleMissingNodeAffinity.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs[i].Spec.Template.Spec.Affinity = nil
	}

	type fields struct {
		bundle             *manifests.Bundle
		supportedPlatforms []platform
	}
	tests := []struct {
		name        string
		fields      fields
		wantWarning bool
		wantError   bool
		warnStrings []string
		errStrings  []string
	}{
		{
			name: "should raise no error or warning when image data matches node affinity",
			fields: fields{
				bundle: validBundle,
				supportedPlatforms: []platform{
					{"linux", "amd64"},
					{"linux", "arm64"},
					{"linux", "ppc64le"},
					{"linux", "s390x"},
				},
			},
		},
		{
			name: "should raise warning when the node affinity information is missing or invalid",
			fields: fields{
				bundle: validBundleMissingNodeAffinity,
				supportedPlatforms: []platform{
					{"linux", "amd64"},
					{"linux", "arm64"},
					{"linux", "ppc64le"},
					{"linux", "s390x"},
				},
			},
			wantWarning: true,
			warnStrings: []string{"check if the CSV has a missing or invalid node affinity configuration for the image: \"quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b\". The image data suggests the following platforms are supported: [\"linux/amd64\" \"linux/arm64\" \"linux/ppc64le\" \"linux/s390x\"]"},
		},
		{
			name: "should raise a warning if node affinity supports more arches than the image data",
			fields: fields{
				bundle: validBundle,
				supportedPlatforms: []platform{
					{"linux", "amd64"},
					{"linux", "arm64"},
				},
			},
			wantWarning: true,
			warnStrings: []string{"the CSV includes [\"linux/ppc64le\" \"linux/s390x\"] in the node affinity configuration for the image: \"quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b\", but the image data suggests the following platforms are supported: [\"linux/amd64\" \"linux/arm64\"]"},
		},
		{
			name: "should raise a warning if the image data supports more arches than the node affinity",
			fields: fields{
				bundle: validBundle,
				supportedPlatforms: []platform{
					{"other", "amd64"},
					{"other", "arm64"},
					{"linux", "amd64"},
					{"linux", "arm64"},
					{"linux", "ppc64le"},
					{"linux", "s390x"},
				},
			},
			wantWarning: true,
			warnStrings: []string{"the image data indicates [\"other/amd64\" \"other/arm64\"] is supported for the image: \"quay.io/coreos/etcd-operator2@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b\", but the node affinity configuration for the image only specifies [\"linux/amd64\" \"linux/arm64\" \"linux/ppc64le\" \"linux/s390x\"]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &multiArchValidator{
				bundle: tt.fields.bundle,
			}

			data.loadImagesFromCSV()

			// Mock inspected platform
			for key, _ := range data.managerImages {
				data.managerImages[key] = tt.fields.supportedPlatforms
			}

			data.checkNodeAffinity(data.managerImages)

			t.Log(data.warns)
			t.Log(data.errors)

			require.Equal(t, tt.wantWarning, len(data.warns) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(data.warns))
				for _, w := range data.warns {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(data.errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(data.errors))
				for _, w := range data.errors {
					wString := w.Error()
					require.Contains(t, tt.errStrings, wString)
				}
			}
		})

	}
}

package validate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BundleValidator struct {
	fileName string
	Manifest Manifest
}

var _ validator.Validator = &BundleValidator{}

func (v *BundleValidator) Validate() (results []validator.ManifestResult) {

	result := bundleInspect(v.Manifest)
	if result.Name == "" {
		result.Name = v.Manifest.Name
	}
	results = append(results, result)

	return results
}

func (v *BundleValidator) AddObjects(objs ...interface{}) validator.Error {
	// TODO: define addObjects for bundle.go
	return validator.Error{}
}

func (v BundleValidator) Name() string {
	return "Bundle Validator"
}

func (v BundleValidator) FileName() string {
	return v.fileName
}

func (v BundleValidator) Unmarshal(rawYaml []byte) (interface{}, error) {
	return nil, fmt.Errorf("Error: unsupported operation; unmarshal not defined for bundle validator")
}

func bundleInspect(manifest Manifest) validator.ManifestResult {
	manifestResult := validator.ManifestResult{}
	csvReplacesMap := make(map[string]string)
	var csvsInBundle []string
	for _, bundle := range manifest.Bundle {
		csv, err := readAndUnmarshalCSV(bundle.CSV)
		if err != (validator.Error{}) {
			manifestResult.Errors = append(manifestResult.Errors, err)
			return manifestResult
		}
		csvsInBundle = append(csvsInBundle, csv.ObjectMeta.Name)
		csvReplacesMap[bundle.CSV] = csv.Spec.Replaces
		if csv.ObjectMeta.Name == csv.Spec.Replaces {
			manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: `spec.replaces` field matches its own `metadata.Name` for %s CSV. It should contain `metadata.Name` of the old CSV to be replaced", bundle.CSV)))
		}
		manifestResult = validateOwnedCRDs(bundle, csv, manifestResult)
	}
	manifestResult = checkReplacesForCSVs(csvReplacesMap, csvsInBundle, manifestResult)
	manifestResult = checkDefaultChannelInBundle(manifest.Package, csvsInBundle, manifestResult)
	return manifestResult
}

func checkDefaultChannelInBundle(pkgName string, csvsInBundle []string, manifestResult validator.ManifestResult) validator.ManifestResult {
	rawYaml, err := ioutil.ReadFile(pkgName)
	if err != nil {
		manifestResult.Errors = append(manifestResult.Errors, validator.IOError(fmt.Sprintf("Error in reading %s file:   #%s ", pkgName, err), pkgName))
		return manifestResult
	}
	v := &PackageValidator{}
	pkg, err := v.Unmarshal(rawYaml)
	if err != nil {
		manifestResult.Errors = append(manifestResult.Errors, validator.InvalidParse(fmt.Sprintf("Error unmarshalling YAML to package manifest type for %s file:  #%s ", pkgName, err), pkgName))
		return manifestResult
	}
	if pkg, ok := pkg.(registry.PackageManifest); ok {
		for _, channel := range pkg.Channels {
			if !isStringPresent(csvsInBundle, channel.CurrentCSVName) {
				manifestResult.Errors = append(manifestResult.Errors, validator.InvalidBundle(fmt.Sprintf("Error: currentCSV `%s` for channel name `%s` in package `%s` not found in manifest", channel.CurrentCSVName, channel.Name, pkg.PackageName), channel.CurrentCSVName))
			}
		}
	}
	return manifestResult
}

func validateOwnedCRDs(bundle ManifestBundle, csv v1alpha1.ClusterServiceVersion, manifestResult validator.ManifestResult) validator.ManifestResult {
	ownedCrdNames := getOwnedCustomResourceDefintionNames(csv)
	bundleCrdNames, err := getBundleCRDNames(bundle)
	if err != (validator.Error{}) {
		manifestResult.Errors = append(manifestResult.Errors, err)
		return manifestResult
	}

	// validating names
	for _, ownedCrd := range ownedCrdNames {
		if !bundleCrdNames[ownedCrd] {
			manifestResult.Errors = append(manifestResult.Errors, validator.InvalidBundle(fmt.Sprintf("Error: owned crd (%s) not found in bundle %s", ownedCrd, bundle.Version), ownedCrd))
		} else {
			delete(bundleCrdNames, ownedCrd)
		}
	}
	// CRDs not defined in the CSV present in the bundle
	if len(bundleCrdNames) != 0 {
		for crd, _ := range bundleCrdNames {
			manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidBundle(fmt.Sprintf("Warning: `%s` crd present in bundle `%s` not defined in csv", crd, bundle.Version), crd))
		}
	}
	return manifestResult
}

type CRDObjectMeta struct {
	metav1.ObjectMeta `json:"metadata"`
}

func getOwnedCustomResourceDefintionNames(csv v1alpha1.ClusterServiceVersion) []string {
	var names []string
	for _, ownedCrd := range csv.Spec.CustomResourceDefinitions.Owned {
		names = append(names, ownedCrd.Name)
	}
	return names
}

func getBundleCRDNames(bundle ManifestBundle) (map[string]bool, validator.Error) {
	bundleCrdNames := make(map[string]bool)
	for _, crdFileName := range bundle.CRDs {
		parsedName := CRDObjectMeta{}
		rawYaml, err := ioutil.ReadFile(crdFileName)
		if err != nil {
			return nil, validator.IOError(fmt.Sprintf("Error in reading %s file:   #%s ", crdFileName, err), crdFileName)
		}
		rawJson, err := yaml.YAMLToJSON(rawYaml)
		if err != nil {
			return nil, validator.InvalidParse(fmt.Sprintf("Error in converting to JSON for %s file:   #%s ", crdFileName, err), crdFileName)
		}

		if err := json.Unmarshal(rawJson, &parsedName); err != nil {
			return nil, validator.InvalidParse(fmt.Sprintf("Error parsing object meta names for %s file:   #%s ", crdFileName, err), crdFileName)
		}
		bundleCrdNames[parsedName.Name] = true
	}
	return bundleCrdNames, validator.Error{}
}

func readAndUnmarshalCSV(pathCSV string) (v1alpha1.ClusterServiceVersion, validator.Error) {
	rawYaml, err := ioutil.ReadFile(pathCSV)
	if err != nil {
		return v1alpha1.ClusterServiceVersion{}, validator.IOError(fmt.Sprintf("Error in reading %s file:   #%s ", pathCSV, err), pathCSV)
	}
	v := &CSVValidator{}
	csv, err := v.Unmarshal(rawYaml)
	if err != nil {
		return v1alpha1.ClusterServiceVersion{}, validator.InvalidParse(fmt.Sprintf("Error unmarshalling YAML to OLM's csv type for %s file:  #%s ", pathCSV, err), pathCSV)
	}
	if csv, ok := csv.(v1alpha1.ClusterServiceVersion); ok {
		return csv, validator.Error{}
	}
	return v1alpha1.ClusterServiceVersion{}, validator.InvalidParse(fmt.Sprintf("Error unmarshalling YAML to OLM's csv type for %s file:  #%s ", pathCSV, err), pathCSV)
}

// checkReplacesForCSVs generates an error if value of the `replaces` field in the
// csv does not match the `metadata.Name` field of the old csv to be replaced.
// It also generates a warning if the `replaces` field of a csv is empty.
func checkReplacesForCSVs(csvReplacesMap map[string]string, csvsInBundle []string, manifestResult validator.ManifestResult) validator.ManifestResult {
	for pathCSV, replaces := range csvReplacesMap {
		if replaces == "" {
			manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: `spec.replaces` field not present in %s csv. If this csv replaces an old version, populate this field with the `metadata.Name` of the old csv", pathCSV)))
		} else {
			if !isStringPresent(csvsInBundle, replaces) {
				manifestResult.Errors = append(manifestResult.Errors, validator.InvalidCSV(fmt.Sprintf("Error: `%s` mentioned in the `spec.replaces` field of %s csv not present in the manifest", replaces, pathCSV)))
			}
		}
	}
	return manifestResult
}

func isStringPresent(list []string, val string) bool {
	for _, str := range list {
		if val == str {
			return true
		}
	}
	return false
}

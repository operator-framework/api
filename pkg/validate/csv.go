package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type CSVValidator struct {
	fileName string
	csvs     []v1alpha1.ClusterServiceVersion
}

var _ validator.Validator = &CSVValidator{}

func (v *CSVValidator) Validate() (results []validator.ManifestResult) {
	for _, csv := range v.csvs {
		result := csvInspect(csv)
		if result.Name == "" {
			result.Name = csv.GetName()
		}
		results = append(results, result)
	}
	return results
}

func (v *CSVValidator) AddObjects(objs ...interface{}) validator.Error {
	for _, o := range objs {
		switch t := o.(type) {
		case v1alpha1.ClusterServiceVersion:
			v.csvs = append(v.csvs, t)
		case *v1alpha1.ClusterServiceVersion:
			v.csvs = append(v.csvs, *t)
		}
	}
	return validator.Error{}
}

func (v CSVValidator) Name() string {
	return "ClusterServiceVersion Validator"
}

func (v CSVValidator) FileName() string {
	return v.fileName
}

func (v CSVValidator) Unmarshal(rawYaml []byte) (interface{}, error) {
	var csv v1alpha1.ClusterServiceVersion

	rawJson, err := yaml.YAMLToJSON(rawYaml)
	if err != nil {
		return v1alpha1.ClusterServiceVersion{}, fmt.Errorf("error parsing raw YAML to Json: %s", err)
	}
	if err := json.Unmarshal(rawJson, &csv); err != nil {
		return v1alpha1.ClusterServiceVersion{}, fmt.Errorf("error parsing CSV (JSON) : %s", err)
	}
	return csv, nil
}

// Iterates over the given CSV. Returns a ManifestResult type object.
func csvInspect(csv v1alpha1.ClusterServiceVersion) validator.ManifestResult {

	// validate example annotations ("alm-examples", "olm.examples").
	manifestResult := validateExamplesAnnotations(csv)

	// validate installModes
	manifestResult = validateInstallModes(csv, manifestResult)

	// check missing optional/mandatory fields.
	fieldValue := reflect.ValueOf(csv)

	switch fieldValue.Kind() {
	case reflect.Struct:
		return checkMissingFields(fieldValue, "", manifestResult)
	default:
		errs := []validator.Error{
			validator.InvalidCSV("Error: input file is not a valid CSV"),
		}

		return validator.ManifestResult{Errors: errs, Warnings: nil}
	}
}

// Recursive function that traverses a nested struct passed in as reflect value, and reports for errors/warnings
// in case of null struct field values.
func checkMissingFields(v reflect.Value, parentStructName string, log validator.ManifestResult) validator.ManifestResult {

	for i := 0; i < v.NumField(); i++ {

		fieldValue := v.Field(i)

		tag := v.Type().Field(i).Tag.Get("json")
		// Ignore fields that are subsets of a primitive field.
		if tag == "" {
			continue
		}

		fields := strings.Split(tag, ",")
		isOptionalField := containsStrict(fields, "omitempty")
		emptyVal := isEmptyValue(fieldValue)

		newParentStructName := ""
		if parentStructName == "" {
			newParentStructName = v.Type().Field(i).Name
		} else {
			newParentStructName = parentStructName + "." + v.Type().Field(i).Name
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			log = updateLog(log, "struct", newParentStructName, emptyVal, isOptionalField)
			if emptyVal {
				continue
			}
			log = checkMissingFields(fieldValue, newParentStructName, log)
		default:
			log = updateLog(log, "field", newParentStructName, emptyVal, isOptionalField)
		}
	}
	return log
}

// Returns updated error log with missing optional/mandatory field/struct objects.
func updateLog(log validator.ManifestResult, typeName string, newParentStructName string, emptyVal bool, isOptionalField bool) validator.ManifestResult {

	if emptyVal && isOptionalField {
		// TODO: update the value field (typeName).
		log.Warnings = append(log.Warnings, validator.OptionalFieldMissing(fmt.Sprintf("Warning: optional %s missing: (%s)", typeName, newParentStructName), newParentStructName, typeName))
	} else if emptyVal && !isOptionalField {
		if newParentStructName != "Status" {
			// TODO: update the value field (typeName).
			log.Errors = append(log.Errors, validator.MandatoryFieldMissing(fmt.Sprintf("Error: mandatory %s missing: (%s)", typeName, newParentStructName), newParentStructName, typeName))
		}
	}
	return log
}

// Takes in a string slice and checks if a string (x) is present in the slice.
// Return true if the string is present in the slice.
func containsStrict(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Uses reflect package to check if the value of the object passed is null, returns a boolean accordingly.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		// Check if the value for 'Spec.InstallStrategy.StrategySpecRaw' field is present. This field is a RawMessage value type. Without a value, the key is explicitly set to 'null'.
		if fieldValue, ok := v.Interface().(json.RawMessage); ok {
			valString := string(fieldValue)
			if valString == "null" {
				return true
			}
		}
		return v.Len() == 0
	// Currently the only CSV field with integer type is containerPort. Operator Verification Library raises a warning if containerPort field is missisng or if its value is 0.
	// It is an optional field so the user can ignore the warning saying this field is missing if they intend to use port 0.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if !isEmptyValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		panic(fmt.Sprintf("%v kind is not supported.", v.Kind()))
	}
}

// validateExamplesAnnotations compares alm/olm example annotations with provided APIs given
// by Spec.CustomResourceDefinitions.Owned and Spec.APIServiceDefinitions.Owned.
func validateExamplesAnnotations(csv v1alpha1.ClusterServiceVersion) (manifestResult validator.ManifestResult) {
	var examples []v1beta1.CustomResourceDefinition
	var annotationsExamples string
	annotations := csv.ObjectMeta.GetAnnotations()
	// Return right away if no examples annotations are found.
	if len(annotations) == 0 {
		manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: example annotations not found for %s csv", csv.GetName())))
		return
	}
	// Expect either `alm-examples` or `olm.examples` but not both
	// If both are present, `alm-examples` will be used
	if value, ok := annotations["alm-examples"]; ok {
		annotationsExamples = value
		if _, ok = annotations["olm.examples"]; ok {
			// both `alm-examples` and `olm.examples` are present
			manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: both `alm-examples` and `olm.examples` are present in %s CSV. Defaulting to `alm-examples` and ignoring `olm.examples`", csv.GetName())))
		}
	} else {
		annotationsExamples = annotations["olm.examples"]
	}

	// Can't find examples annotations, simply return
	if annotationsExamples == "" {
		manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: example annotations not found for %s csv", csv.GetName())))
		return
	}

	if err := json.Unmarshal([]byte(annotationsExamples), &examples); err != nil {
		manifestResult = getManifestResult(validator.InvalidParse(fmt.Sprintf("Error: parsing example annotations to %T type:  %s ", examples, err), nil))
		return
	}

	providedAPIs, manRes := getProvidedAPIs(csv, manifestResult)

	parsedExamples, manRes := parseExamplesAnnotations(examples, manifestResult)
	if len(manRes.Errors) != 0 || len(manRes.Warnings) != 0 {
		return manRes
	}

	return matchGVKProvidedAPIs(parsedExamples, providedAPIs, manifestResult)
}

func getProvidedAPIs(csv v1alpha1.ClusterServiceVersion, manifestResult validator.ManifestResult) (map[schema.GroupVersionKind]struct{}, validator.ManifestResult) {
	provided := map[schema.GroupVersionKind]struct{}{}

	for _, owned := range csv.Spec.CustomResourceDefinitions.Owned {
		parts := strings.SplitN(owned.Name, ".", 2)
		if len(parts) < 2 {
			manifestResult.Errors = append(manifestResult.Errors, validator.InvalidParse(fmt.Sprintf("Error: couldn't parse plural.group from crd name: %s", owned.Name), owned.Name))
			continue
		}
		provided[schema.GroupVersionKind{Group: parts[1], Version: owned.Version, Kind: owned.Kind}] = struct{}{}
	}

	for _, api := range csv.Spec.APIServiceDefinitions.Owned {
		provided[schema.GroupVersionKind{Group: api.Group, Version: api.Version, Kind: api.Kind}] = struct{}{}
	}

	return provided, manifestResult
}

func parseExamplesAnnotations(examples []v1beta1.CustomResourceDefinition, manifestResult validator.ManifestResult) (map[schema.GroupVersionKind]struct{}, validator.ManifestResult) {
	parsed := map[schema.GroupVersionKind]struct{}{}
	for _, value := range examples {
		parts := strings.SplitN(value.APIVersion, "/", 2)
		if len(parts) < 2 {
			manifestResult.Errors = append(manifestResult.Errors, validator.InvalidParse(fmt.Sprintf("Error: couldn't parse group/version from crd kind: %s", value.Kind), value.Kind))
			continue
		}
		parsed[schema.GroupVersionKind{Group: parts[0], Version: parts[1], Kind: value.Kind}] = struct{}{}
	}

	return parsed, manifestResult
}

func matchGVKProvidedAPIs(examples map[schema.GroupVersionKind]struct{}, providedAPIs map[schema.GroupVersionKind]struct{}, manifestResult validator.ManifestResult) validator.ManifestResult {
	for key := range examples {
		if _, ok := providedAPIs[key]; !ok {
			manifestResult.Errors = append(manifestResult.Errors, validator.InvalidOperation(fmt.Sprintf("Error: couldn't match %v in provided APIs list: %v", key, providedAPIs), key))
			continue
		}
	}
	return manifestResult
}

func getManifestResult(errs ...validator.Error) validator.ManifestResult {
	errList := append([]validator.Error{}, errs...)
	return validator.ManifestResult{Errors: errList, Warnings: nil}
}

func validateInstallModes(csv v1alpha1.ClusterServiceVersion, manifestResult validator.ManifestResult) validator.ManifestResult {
	// var installModeSet v1alpha1.InstallModeSet
	installModeSet := make(v1alpha1.InstallModeSet)
	for _, installMode := range csv.Spec.InstallModes {
		if _, ok := installModeSet[installMode.Type]; ok {
			manifestResult.Errors = append(manifestResult.Errors, validator.InvalidCSV(fmt.Sprintf("Error: duplicate install modes present in %s csv", csv.GetName())))
		} else {
			installModeSet[installMode.Type] = installMode.Supported
		}
	}

	// installModes not found, return with a warning
	if len(installModeSet) == 0 {
		manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidCSV(fmt.Sprintf("Warning: install modes not found for %s csv", csv.GetName())))
		return manifestResult
	}

	// all installModes should not be `false`
	if checkAllFalseForInstallModeSet(installModeSet) {
		manifestResult.Errors = append(manifestResult.Errors, validator.InvalidCSV(fmt.Sprintf("Error: none of InstallModeTypes are supported for %s csv", csv.GetName())))
	}
	return manifestResult
}

func checkAllFalseForInstallModeSet(installModeSet v1alpha1.InstallModeSet) bool {
	for _, isSupported := range installModeSet {
		if isSupported {
			return false
		}
	}
	return true
}

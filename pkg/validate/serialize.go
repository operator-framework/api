package validate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

// Represents verification result for each of the yaml files from the manifest bundle.
type manifestResult struct {
	errors   []missingTypeError
	warnings []missingTypeError
}

// Represents a warning or an error in a yaml file.
type missingTypeError struct {
	err         string
	typeName    string
	path        string
	isMandatory bool
}

// ValidateCSVManifest takes in name of the yaml file to be validated, reads
// it, and calls the unmarshal function on rawYaml.
func ValidateCSVManifest(yamlFileName string) error {
	rawYaml, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return fmt.Errorf("Error in reading %s file:   #%s ", yamlFileName, err)
	}

	// Value returned is a marshaled go type (CSV Struct).
	csv, err := unmarshal(rawYaml)
	if err != nil {
		return fmt.Errorf("Error unmarshalling YAML to OLM's csv type for %s file:  #%s ", yamlFileName, err)
	}

	v := &CSVValidator{}
	if err = v.AddObjects(csv); err != nil {
		return err
	}
	fmt.Println("Running", v.Name())
	if err = v.Validate(); err != nil {
		return err
	}
	fmt.Printf("%s is verified.", yamlFileName)
	return nil
}

// Iterates over the list of warnings and errors.
func getErrorsFromManifestResult(err []missingTypeError) {
	for _, v := range err {
		assertTypeToGetValue(v)
	}
}

// Asserts type to get the underlying field value.
func assertTypeToGetValue(v interface{}) {
	if v, ok := v.(missingTypeError); ok {
		fmt.Println(v)
	}
}

// Unmarshal takes in a raw YAML file and deserializes it to OLM's ClusterServiceVersion type.
// Throws an error if:
// (1) the yaml file can not be converted to json.
// (2) there is a problem while unmarshalling in go type.
// Returns an object of type olm.ClusterServiceVersion.
func unmarshal(rawYAML []byte) (olm.ClusterServiceVersion, error) {

	var csv olm.ClusterServiceVersion

	rawJson, err := yaml.YAMLToJSON(rawYAML)
	if err != nil {
		fmt.Printf("error parsing raw YAML to Json: %s", err)
		return csv, err
	}
	if err := json.Unmarshal(rawJson, &csv); err != nil {
		return csv, fmt.Errorf("error parsing CSV list (JSON) : %s", err)
	}

	return csv, nil
}

// missingTypeError strut implements the Error interface to define custom error formatting.
func (err missingTypeError) Error() string {
	if err.isMandatory {
		return fmt.Sprintf("Error: Mandatory %s Missing (%s)", err.typeName, err.path)
	} else {
		return fmt.Sprintf("Warning: Optional %s Missing (%s)", err.typeName, err.path)
	}
}

// Iterates over the given CSV. Returns a manifestResult type object.
func csvInspect(val interface{}) manifestResult {

	fieldValue := reflect.ValueOf(val)

	switch fieldValue.Kind() {
	case reflect.Struct:
		return checkMissingFields(fieldValue, "", manifestResult{})
	default:
		err := []missingTypeError{{"Error: input file is not a valid CSV.", "", "", false}}
		return manifestResult{errors: err, warnings: nil}
	}
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

// Recursive function that traverses a nested struct passed in as reflect value, and reports for errors/warnings
// in case of null struct field values.
// Returns a log of errors as slice of strings.
func checkMissingFields(v reflect.Value, parentStructName string, log manifestResult) manifestResult {

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
			log = updateLog(log, "Struct", newParentStructName, emptyVal, isOptionalField)
			if emptyVal {
				continue
			}
			log = checkMissingFields(fieldValue, newParentStructName, log)
		default:
			log = updateLog(log, "Field", newParentStructName, emptyVal, isOptionalField)
		}
	}
	return log
}

// Returns updated error log with missing optional/mandatory field/struct objects.
func updateLog(log manifestResult, typeName string, newParentStructName string, emptyVal bool, isOptionalField bool) manifestResult {

	if emptyVal && isOptionalField {
		errString := fmt.Sprintf("Warning: Optional %s Missing.", typeName)
		log.warnings = append(log.warnings, missingTypeError{errString, typeName, newParentStructName, false})
	} else if emptyVal && !isOptionalField {
		if newParentStructName != "Status" {
			errString := fmt.Sprintf("Error: Mandatory %s Missing.", typeName)
			log.errors = append(log.errors, missingTypeError{errString, typeName, newParentStructName, true})
		}
	}
	return log
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

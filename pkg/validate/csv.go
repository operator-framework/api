package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/pkg/errors"

	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

type CSVValidator struct {
	csvs []olm.ClusterServiceVersion
}

var _ validator.Validator = &CSVValidator{}

func (v *CSVValidator) Validate() (results []validator.ManifestResult) {
	for _, csv := range v.csvs {
		// Contains error logs for all missing optional and mandatory fields.
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
		case olm.ClusterServiceVersion:
			v.csvs = append(v.csvs, t)
		case *olm.ClusterServiceVersion:
			v.csvs = append(v.csvs, *t)
		}
	}
	return validator.Error{}
}

func (v CSVValidator) Name() string {
	return "ClusterServiceVersion Validator"
}

// Iterates over the given CSV. Returns a ManifestResult type object.
func csvInspect(val interface{}) validator.ManifestResult {

	fieldValue := reflect.ValueOf(val)

	switch fieldValue.Kind() {
	case reflect.Struct:
		return checkMissingFields(fieldValue, "", validator.ManifestResult{})
	default:
		errs := []validator.Error{
			validator.InvalidCSV("Error: input file is not a valid CSV."),
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
func updateLog(log validator.ManifestResult, typeName string, newParentStructName string, emptyVal bool, isOptionalField bool) validator.ManifestResult {

	if emptyVal && isOptionalField {
		err := errors.Errorf("Warning: Optional %s Missing", typeName)
		// TODO: update the value field (typeName).
		log.Warnings = append(log.Warnings, validator.OptionalFieldMissing(newParentStructName, typeName, err.Error()))
	} else if emptyVal && !isOptionalField {
		if newParentStructName != "Status" {
			err := errors.Errorf("Error: Mandatory %s Missing", typeName)
			// TODO: update the value field (typeName).
			log.Errors = append(log.Errors, validator.MandatoryFieldMissing(newParentStructName, typeName, err.Error()))
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

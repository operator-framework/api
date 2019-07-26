package validate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"

	"github.com/ghodss/yaml"
	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

// ValidateCSVManifest takes in name of the yaml file to be validated, reads
// it, and calls the unmarshal function on rawYaml.
func ValidateCSVManifest(yamlFileName string) validator.Error {
	rawYaml, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return validator.IOError(fmt.Sprintf("Error in reading %s file:   #%s ", yamlFileName, err), yamlFileName)
	}

	// Value returned is a marshaled go type (CSV Struct).
	csv, err := unmarshal(rawYaml)
	if err != nil {
		return validator.InvalidParse(fmt.Sprintf("Error unmarshalling YAML to OLM's csv type for %s file:  #%s ", yamlFileName, err), yamlFileName)
	}

	v := &CSVValidator{}
	if err := v.AddObjects(csv); err != (validator.Error{}) {
		return err // TODO: update when 'AddObjects' returns an actual error.
	}
	fmt.Println("Running", v.Name())
	for _, errorLog := range v.Validate() {
		fmt.Println("Validating CSV", errorLog.Name)

		getErrorsFromManifestResult(errorLog.Warnings)

		// There is no mandatory field thats missing if errorLog.errors is nil.
		if errorLog.Errors != nil {
			fmt.Println()
			getErrorsFromManifestResult(errorLog.Errors)
			fmt.Printf("Populate all the mandatory fields missing from CSV %s.", csv.GetName())
			return validator.Error{}
		}
	}
	fmt.Printf("%s is verified.\n", yamlFileName)
	return validator.Error{}
}

// Iterates over the list of warnings and errors.
func getErrorsFromManifestResult(err []validator.Error) {
	for _, v := range err {
		assertTypeToGetValue(v)
	}
}

// Asserts type to get the underlying field value.
func assertTypeToGetValue(v interface{}) {
	if v, ok := v.(validator.Error); ok {
		fmt.Println(v.String())
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
		return olm.ClusterServiceVersion{}, err
	}
	if err := json.Unmarshal(rawJson, &csv); err != nil {
		return olm.ClusterServiceVersion{}, fmt.Errorf("error parsing CSV list (JSON) : %s", err)
	}
	return csv, nil
}

package validate

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

// Verify takes in name of the yaml file to be validated, reads it, and calls the unmarshal function on rawYaml.
func Verify(yamlFileName string) error {
	rawYaml, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return fmt.Errorf("Error in reading %s file:   #%s ", yamlFileName, err)
	}

	// Value returned is a marshaled go type (CSV Struct).
	_, err = unmarshal(rawYaml)
	if err != nil {
		return fmt.Errorf("Error unmarshalling YAML to OLM's csv type for %s file:  #%s ", yamlFileName, err)
	}

	fmt.Printf("%s is verified", yamlFileName)
	return nil
}

// Unmarshal takes in a raw YAML file and deserializes it to OLM's ClusterServiceVersion type.
// Throws an error if:
// (1) the yaml file can not be converted to json.
// (2) there is a problem while unmarshalling in go type.
// Returns an object of type olm.ClusterServiceVersion.
func unmarshal(rawYAML []byte) (*olm.ClusterServiceVersion, error) {

	var csv olm.ClusterServiceVersion

	if err := yaml.Unmarshal(rawYAML, &csv); err != nil {
		return nil, fmt.Errorf("error parsing CSV list (JSON) : %s", err)
	}

	return &csv, nil

}

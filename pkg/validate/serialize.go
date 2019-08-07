package validate

import (
	"fmt"
	"io/ioutil"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
)

func Validate(v validator.Validator) (manifestResult validator.ManifestResult) {
	fmt.Printf("\nRunning %s\n", v.Name())
	fmt.Printf("Validating %s\n\n", v.FileName())
	rawYaml, err := ioutil.ReadFile(v.FileName())
	if err != nil {
		manifestResult.Errors = append(manifestResult.Errors, validator.IOError(fmt.Sprintf("Error in reading %s file:   #%s ", v.FileName(), err), v.FileName()))
		getErrorsFromManifestResult(manifestResult.Errors)
		return
	}

	// Value returned is a marshaled go type.
	unmarshalledObject, err := v.Unmarshal(rawYaml)
	if err != nil {
		manifestResult.Errors = append(manifestResult.Errors, validator.InvalidParse(fmt.Sprintf("Error unmarshalling YAML for %s file:  #%s ", v.FileName(), err), v.FileName()))
		getErrorsFromManifestResult(manifestResult.Errors)
		return
	}

	if err := v.AddObjects(unmarshalledObject); err != (validator.Error{}) {
		manifestResult.Errors = append(manifestResult.Errors, err)
		getErrorsFromManifestResult(manifestResult.Errors)
		return // TODO: update when 'AddObjects' returns an actual error.
	}

	for _, errorLog := range v.Validate() {

		getErrorsFromManifestResult(errorLog.Warnings)

		if len(errorLog.Errors) != 0 {
			fmt.Println()
			getErrorsFromManifestResult(errorLog.Errors)
		} else {
			fmt.Printf("\n%s is verified\n", v.FileName())
		}
	}
	return
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

func validateBundle(manifest Manifest) []validator.ManifestResult {
	v := &BundleValidator{Manifest: manifest}
	manifestResult := v.Validate()
	for _, errorLog := range manifestResult {
		fmt.Printf("\nValidating `%s` Manifest\n", errorLog.Name)
		fmt.Println()
		if len(errorLog.Warnings) != 0 {
			getErrorsFromManifestResult(errorLog.Warnings)
		}

		if len(errorLog.Errors) != 0 {
			fmt.Println()
			getErrorsFromManifestResult(errorLog.Errors)
			fmt.Printf("Invalid manifest: `%s`\n", errorLog.Name)
		} else {
			fmt.Printf("`%s` manifest verified", errorLog.Name)
		}
	}
	return manifestResult
}

func parseManifestDirectory(manifestDirectory string) (Manifest, []validator.ManifestResult) {
	manifestResultList := []validator.ManifestResult{}
	fmt.Printf("Parsing `%s` operator manifest\n\n", manifestDirectory)
	manifest, manifestResultFromDirectoryParse := ParseDir(manifestDirectory)

	if len(manifestResultFromDirectoryParse.Errors) != 0 || len(manifestResultFromDirectoryParse.Warnings) != 0 {
		manifestResultList = append(manifestResultList, manifestResultFromDirectoryParse)
		getErrorsFromManifestResult(manifestResultFromDirectoryParse.Warnings)
		if len(manifestResultFromDirectoryParse.Errors) != 0 {
			getErrorsFromManifestResult(manifestResultFromDirectoryParse.Errors)
			fmt.Printf("Invalid operator manifest structure for `%s`\n", manifestDirectory)
			return Manifest{}, manifestResultList
		}
	}
	return manifest, manifestResultList
}

func ValidateManifest(manifestDirectory string) []validator.ManifestResult {
	// parse manifest directory
	manifest, manifestResultList := parseManifestDirectory(manifestDirectory)
	for _, manifestResult := range manifestResultList {
		if len(manifestResult.Errors) != 0 {
			return manifestResultList
		}
	}

	var result []validator.ManifestResult
	// validate individual bundle files
	for _, bundle := range manifest.Bundle {
		validators := []validator.Validator{&CSVValidator{fileName: bundle.CSV}}
		for _, crd := range bundle.CRDs {
			validators = append(validators, &CRDValidator{fileName: crd})
		}
		for _, validator := range validators {
			result = append(result, Validate(validator))
		}
	}
	var pkgValidator validator.Validator
	pkgValidator = &PackageValidator{fileName: manifest.Package}
	Validate(pkgValidator)

	// validate bundle
	validateBundle(manifest)
	return []validator.ManifestResult{}
}

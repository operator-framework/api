package internal

import (
	goerrors "errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation/errors"
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
)

// GoodPracticesValidator validates the bundle against criteria and suggestions defined as
// good practices for bundles under the operator-framework solutions. (You might give a
// look at https://sdk.operatorframework.io/docs/best-practices/)
//
// This validator will raise an WARNING when:
//
// - The bundle name (CSV.metadata.name) does not follow the naming convention: <operator-name>.v<semver> e.g. memcached-operator.v0.0.1
//
// NOTE: The bundle name must be 63 characters or less because it will be used as k8s ownerref label which only allows max of 63 characters.

var GoodPracticesValidator interfaces.Validator = interfaces.ValidatorFunc(goodPracticesValidator)

func goodPracticesValidator(objs ...interface{}) (results []errors.ManifestResult) {
	for _, obj := range objs {
		switch v := obj.(type) {
		case *manifests.Bundle:
			results = append(results, validateGoodPracticesFrom(v))
		}
	}
	return results
}

func validateGoodPracticesFrom(bundle *manifests.Bundle) errors.ManifestResult {
	result := errors.ManifestResult{}
	if bundle == nil {
		result.Add(errors.ErrInvalidBundle("Bundle is nil", nil))
		return result
	}

	result.Name = bundle.Name

	if bundle.CSV == nil {
		result.Add(errors.ErrInvalidBundle("Bundle csv is nil", bundle.Name))
		return result
	}

	checks := CSVChecks{csv: *bundle.CSV, errs: []error{}, warns: []error{}}

	checks = validateResourceRequests(checks)
	checks = checkBundleName(checks)

	for _, err := range checks.errs {
		result.Add(errors.ErrFailedValidation(err.Error(), bundle.CSV.GetName()))
	}
	for _, warn := range checks.warns {
		result.Add(errors.WarnFailedValidation(warn.Error(), bundle.CSV.GetName()))
	}

	return result
}

// validateResourceRequests will return a WARN when the resource request is not set
func validateResourceRequests(checks CSVChecks) CSVChecks {
	if checks.csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs == nil {
		checks.errs = append(checks.errs, goerrors.New("unable to find a deployment to install in the CSV"))
		return checks
	}
	deploymentSpec := checks.csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs

	for _, dSpec := range deploymentSpec {
		for _, c := range dSpec.Spec.Template.Spec.Containers {
			if c.Resources.Requests == nil || !(len(c.Resources.Requests.Cpu().String()) != 0 && len(c.Resources.Requests.Memory().String()) != 0) {
				msg := fmt.Errorf("unable to find the resource requests for the container: (%s). It is recommended "+
					"to ensure the resource request for CPU and Memory. Be aware that for some clusters configurations "+
					"it is required to specify requests or limits for those values. Otherwise, the system or quota may "+
					"reject Pod creation. More info: https://master.sdk.operatorframework.io/docs/best-practices/managing-resources/", c.Name)
				checks.warns = append(checks.warns, msg)
			}
		}
	}
	return checks
}

// checkBundleName will validate the operator bundle name informed via CSV.metadata.name.
// The motivation for the following check is to ensure that operators authors knows that operator bundles names should
// follow a name and versioning convention
func checkBundleName(checks CSVChecks) CSVChecks {

	// Check if is following the semver
	re := regexp.MustCompile("([0-9]+)\\.([0-9]+)\\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\\.[0-9A-Za-z-]+)*))?(?:\\+[0-9A-Za-z-]+)?$")
	match := re.FindStringSubmatch(checks.csv.Name)

	if len(match) > 0 {
		if _, err := semver.Parse(match[0]); err != nil {
			checks.warns = append(checks.warns, fmt.Errorf("csv.metadata.Name %v is not following the versioning "+
				"convention (MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/", checks.csv.Name))
		}
	} else {
		checks.warns = append(checks.warns, fmt.Errorf("csv.metadata.Name %v is not following the versioning "+
			"convention (MAJOR.MINOR.PATCH e.g 0.0.1): https://semver.org/", checks.csv.Name))
	}

	// Check if its following the name convention
	if len(strings.Split(checks.csv.Name, ".v")) != 2 {
		checks.warns = append(checks.errs, fmt.Errorf("csv.metadata.Name %v is not following the recommended "+
			"naming convention: <operator-name>.v<semver> e.g. memcached-operator.v0.0.1", checks.csv.Name))
	}

	return checks
}

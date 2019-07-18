package validate

import (
	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"

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

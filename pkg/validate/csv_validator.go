package validate

import (
	"fmt"

	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
)

type CSVValidator struct {
	csvs []olm.ClusterServiceVersion
}

var _ Validator = &CSVValidator{}

func (v *CSVValidator) Validate() error {
	for _, csv := range v.csvs {
		// Contains error logs for all missing optional and mandatory fields.
		errorLog := csvInspect(csv)
		getErrorsFromManifestResult(errorLog.warnings)

		// There is no mandatory field thats missing if errorLog.errors is nil.
		if errorLog.errors != nil {
			fmt.Println()
			getErrorsFromManifestResult(errorLog.errors)
			return fmt.Errorf("Populate all the mandatory fields missing from CSV %s.", csv.GetName())
		}
	}
	return nil
}

func (v *CSVValidator) AddObjects(objs ...interface{}) error {
	for _, o := range objs {
		switch t := o.(type) {
		case olm.ClusterServiceVersion:
			v.csvs = append(v.csvs, t)
		case *olm.ClusterServiceVersion:
			v.csvs = append(v.csvs, *t)
		}
	}
	return nil
}

func (v CSVValidator) Name() string {
	return "ClusterServiceVersion Validator"
}

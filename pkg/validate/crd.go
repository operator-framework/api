package validate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/ghodss/yaml"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

type CRDValidator struct {
	fileName string
	crds     []v1beta1.CustomResourceDefinition
}

var _ validator.Validator = &CRDValidator{}

func (v *CRDValidator) Validate() (results []validator.ManifestResult) {
	for _, crd := range v.crds {
		scheme := runtime.NewScheme()
		result := crdInspect(crd, scheme)
		if result.Name == "" {
			result.Name = crd.GetName()
		}
		results = append(results, result)
	}
	return results
}

func (v *CRDValidator) AddObjects(objs ...interface{}) validator.Error {
	for _, o := range objs {
		switch t := o.(type) {
		case v1beta1.CustomResourceDefinition:
			v.crds = append(v.crds, t)
		case *v1beta1.CustomResourceDefinition:
			v.crds = append(v.crds, *t)
		}
	}
	return validator.Error{}
}

func (v CRDValidator) Name() string {
	return "CustomResourceDefinition Validator"
}

func (v CRDValidator) FileName() string {
	return v.fileName
}

func (v CRDValidator) Unmarshal(rawYaml []byte) (interface{}, error) {
	var crd v1beta1.CustomResourceDefinition

	rawJson, err := yaml.YAMLToJSON(rawYaml)
	if err != nil {
		return v1beta1.CustomResourceDefinition{}, fmt.Errorf("error parsing raw YAML to Json: %s", err)
	}
	if err := json.Unmarshal(rawJson, &crd); err != nil {
		return v1beta1.CustomResourceDefinition{}, fmt.Errorf("error parsing CRD (JSON) : %s", err)
	}
	return crd, nil
}

func crdInspect(crd v1beta1.CustomResourceDefinition, scheme *runtime.Scheme) (manifestResult validator.ManifestResult) {
	err := apiextensions.AddToScheme(scheme)
	if err != nil {
		return
	}
	err = v1beta1.AddToScheme(scheme)
	if err != nil {
		return
	}
	unversionedCRD := apiextensions.CustomResourceDefinition{}
	scheme.Converter().Convert(&crd, &unversionedCRD, conversion.SourceToDest, nil)
	errList := validation.ValidateCustomResourceDefinition(&unversionedCRD)
	for _, err := range errList {
		if !strings.Contains(err.Field, "openAPIV3Schema") && !strings.Contains(err.Field, "status") {
			er := validator.Error{Type: validator.ErrorType(err.Type), Field: err.Field, BadValue: err.BadValue, Detail: err.Error()}
			manifestResult.Errors = append(manifestResult.Errors, er)
		}
	}
	return
}

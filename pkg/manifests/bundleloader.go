package manifests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/yaml"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
)

// bundleLoader loads a bundle directory from disk
type bundleLoader struct {
	dir    string
	bundle *Bundle
}

func NewBundleLoader(dir string) bundleLoader {
	return bundleLoader{
		dir: dir,
	}
}

func (b *bundleLoader) LoadBundle() error {
	errs := make([]error, 0)
	if err := filepath.Walk(b.dir, collectWalkErrs(b.LoadBundleWalkFunc, &errs)); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// collectWalkErrs calls the given walk func and appends any non-nil, non skip dir error returned to the given errors slice.
func collectWalkErrs(walk filepath.WalkFunc, errs *[]error) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) (walkErr error) {
		if walkErr = walk(path, f, err); walkErr != nil && walkErr != filepath.SkipDir {
			*errs = append(*errs, walkErr)
			return nil
		}

		return walkErr
	}
}

func (b *bundleLoader) LoadBundleWalkFunc(path string, f os.FileInfo, err error) error {
	if f == nil {
		return fmt.Errorf("invalid file: %v", f)
	}

	if f.IsDir() {
		if strings.HasPrefix(f.Name(), ".") {
			return filepath.SkipDir
		}
		return nil
	}

	if strings.HasPrefix(f.Name(), ".") {
		return nil
	}

	fileReader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to load file %s: %s", path, err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
	csv := unstructured.Unstructured{}

	if err = decoder.Decode(&csv); err != nil {
		return nil
	}

	if csv.GetKind() != operatorsv1alpha1.ClusterServiceVersionKind {
		return nil
	}

	var errs []error
	bundle, err := loadBundle(csv.GetName(), filepath.Dir(path))
	if err != nil {
		errs = append(errs, fmt.Errorf("error loading objs in directory: %s", err))
	}

	if bundle == nil || bundle.CSV == nil {
		errs = append(errs, fmt.Errorf("no bundle csv found"))
		return utilerrors.NewAggregate(errs)
	}

	b.bundle = bundle

	return utilerrors.NewAggregate(errs)
}

// loadBundle takes the directory that a CSV is in and assumes the rest of the objects in that directory
// are part of the bundle.
func loadBundle(csvName string, dir string) (*Bundle, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var errs []error
	bundle := &Bundle{
		Name: csvName,
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		path := filepath.Join(dir, f.Name())
		fileReader, err := os.Open(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to load file %s: %s", path, err))
			continue
		}

		decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
		obj := &unstructured.Unstructured{}
		if err = decoder.Decode(obj); err != nil {
			continue
		}

		bundle.Objects = append(bundle.Objects, obj)

		switch kind := obj.GetKind(); kind {
		case "ClusterServiceVersion":
			if bundle.CSV != nil {
				return nil, fmt.Errorf("invalid bundle: contains multiple CSVs")
			}
			csv := operatorsv1alpha1.ClusterServiceVersion{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(
				obj.UnstructuredContent(),
				&csv)
			if err != nil {
				return nil, fmt.Errorf("unable to parse CSV %s: %s", f.Name(), err.Error())
			}
			bundle.CSV = &csv
		case "CustomResourceDefinition":
			version := obj.GetAPIVersion()
			if version == apiextensionsv1beta1.SchemeGroupVersion.String() {
				crd := apiextensionsv1beta1.CustomResourceDefinition{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(
					obj.UnstructuredContent(),
					&crd)
				if err != nil {
					return nil, fmt.Errorf("unable to parse CRD %s: %s", f.Name(), err.Error())
				}
				bundle.V1beta1CRDs = append(bundle.V1beta1CRDs, &crd)
			} else if version == apiextensionsv1.SchemeGroupVersion.String() {
				crd := apiextensionsv1.CustomResourceDefinition{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(
					obj.UnstructuredContent(),
					&crd)
				if err != nil {
					return nil, fmt.Errorf("unable to parse CRD %s: %s", f.Name(), err.Error())
				}
				bundle.V1CRDs = append(bundle.V1CRDs, &crd)
			} else {
				return nil, fmt.Errorf("unsupported CRD version %s for %s", version, f.Name())
			}
		}
	}

	return bundle, utilerrors.NewAggregate(errs)
}

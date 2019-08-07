package validate

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate/validator"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	yamlForUnmarshalStrict "sigs.k8s.io/yaml"
)

// Manifest represents files in the operator manifest.
type Manifest struct {
	Name string
	// Package stores the name (path) of package yaml file.
	Package string
	// Bundle represents a directory of files with one `ClusterServiceVersion`.
	Bundle map[string]ManifestBundle
}

type ManifestBundle struct {
	// Version stores the CSV version for the bundle.
	Version string
	// List of CustomResourceDefinition file names inside the bundle.
	CRDs []string
	// CSV file name in the bundle.
	CSV string
}

// getFileType identifies the file type and returns it as a string.
func getFileType(filePath string) (string, error) {

	rawYaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error in reading %s file", filePath)
	}
	pkg := registry.PackageManifest{}

	if checkFileTypeWithUnmarshalStrict(rawYaml, &pkg) {
		return "Package", nil
	}

	kind, err := getKindFromFileBytes(rawYaml, filePath)
	if err != nil {
		return "", err
	}
	return kind, nil
}

func getKindFromFileBytes(rawYaml []byte, filePath string) (string, error) {
	u := unstructured.Unstructured{}
	r := bytes.NewReader(rawYaml)
	dec := yaml.NewYAMLOrJSONDecoder(r, 8)
	// There is only one YAML doc if there are no more bytes to be read or EOF
	// is hit.
	if err := dec.Decode(&u); err == nil && r.Len() != 0 {
		return "", fmt.Errorf("error getting TypeMeta from bytes: more than one manifest in bytes")
	} else if err != nil && err != io.EOF {
		return "", fmt.Errorf("error getting TypeMeta from bytes")
	}
	return u.GetKind(), nil
}

func checkFileTypeWithUnmarshalStrict(rawYaml []byte, obj interface{}) bool {
	if err := yamlForUnmarshalStrict.UnmarshalStrict(rawYaml, &obj); err != nil {
		return false
	}
	return true
}

// ParseDir walks through the operator manifest directory, checks its format,
// and populates the Manifest object with relevant file names.
func ParseDir(manifestDirectory string) (Manifest, validator.ManifestResult) {

	countPkg := 0
	manifest := Manifest{}
	manifest.Name = manifestDirectory
	visitedFile := map[string]struct{}{}
	manifestResult := validator.ManifestResult{}
	isManifestResultNameSet := false
	manifest.Bundle = make(map[string]ManifestBundle)
	// parse manifest directory structure
	err := filepath.Walk(manifestDirectory, func(path string, f os.FileInfo, err error) error {

		// set manifest name
		if !isManifestResultNameSet {
			manifestResult.Name = filepath.Base(path)
			isManifestResultNameSet = true
		}
		// create a manifest bundle for each version in the manifest
		if f.IsDir() && path != manifestDirectory {
			if _, ok := manifest.Bundle[path]; !ok {
				bundle := ManifestBundle{}
				bundle.Version = f.Name()
				manifest.Bundle[path] = bundle
			}
		} else if !f.IsDir() {
			fileType, err := getFileType(path)
			if err != nil {
				updateErr := fmt.Sprintf("Error: %s file may not be of ClusterServiceVersion, CustomResourceDefinition, or Package yaml type. If it is supposed to be ClusterServiceVersion or CustomResourceDefinition type, make sure the TypeMeta is correctly defined. If this is a package yaml, instead, make sure it follows the PackageManifest type definition", path)
				manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(updateErr))
				return nil
			}

			directoryPath := filepath.Dir(path)
			switch fileType {
			case "ClusterServiceVersion":
				if _, ok := visitedFile[directoryPath]; ok {
					manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(fmt.Sprintf("Error: more than one CSV in the bundle found at %s bundle", directoryPath)))
					return nil
				} else {
					visitedFile[directoryPath] = struct{}{}
				}
				if bundleObj, ok := manifest.Bundle[directoryPath]; ok {
					bundleObj.CSV = path
					manifest.Bundle[directoryPath] = bundleObj
				} else {
					manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(fmt.Sprintf("Error: %s file at %s path does not align with the operator manifest format", f.Name(), path)))
					return nil
				}
			case "CustomResourceDefinition":
				if bundleObj, ok := manifest.Bundle[directoryPath]; ok {
					bundleObj.CRDs = append(bundleObj.CRDs, path)
					manifest.Bundle[directoryPath] = bundleObj
				} else {
					manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(fmt.Sprintf("Error: %s file at %s path does not align with the operator manifest format", f.Name(), path)))
					return nil
				}
			case "Package":
				countPkg++
				if countPkg > 1 {
					manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(fmt.Sprintf("Error: more than one package yaml file in the manifest; found at %s", path)))
					return nil
				}
				manifest.Package = path
				return nil
			default:
				// return a `Warning` for files other than CSV, CRD, package yaml present in the manifest. If required, we can keep a list of these file
				// paths and remove them from the manifest.
				manifestResult.Warnings = append(manifestResult.Warnings, validator.InvalidManifestStructure(fmt.Sprintf("Warning: %s file at %s path is not a ClusterServiceVersion, CustomResourceDefinition, or Package yaml type", f.Name(), path)))
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	if countPkg == 0 {
		manifestResult.Errors = append(manifestResult.Errors, validator.InvalidManifestStructure(fmt.Sprintf("Error: no package yaml in `%s` manifest", manifestDirectory)))
	}
	return manifest, manifestResult
}

package metadata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

// The bundle root-relative default annotations.yaml path.
var defaultRelPath = filepath.Join("metadata", "annotations.yaml")

// FindAnnotations walks bundleRoot searching for bundle metadata in an annotations.yaml file,
// and returns this metadata and its path if found. If one is not found, an error is returned.
func FindAnnotations(bundleRoot string) (map[string]string, string, error) {
	return findAnnotations(afero.NewOsFs(), bundleRoot)
}

func findAnnotations(fs afero.Fs, bundleRoot string) (map[string]string, string, error) {
	// Check the default path first, and return annotations if they were found or an error if that error
	// is not because the path does not exist (it exists or there was an unmarshalling error).
	annotationsPath := filepath.Join(bundleRoot, defaultRelPath)
	annotations, err := readAnnotations(fs, annotationsPath)
	if (err == nil && len(annotations) != 0) || (err != nil && !errors.Is(err, os.ErrNotExist)) {
		return annotations, annotationsPath, err
	}

	// Annotations are not at the default path, so search recursively.
	annotations = make(map[string]string)
	annotationsPath = ""
	err = afero.Walk(fs, bundleRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories and hidden files, or if annotations were already found.
		if len(annotations) != 0 || info.IsDir() || strings.HasPrefix(path, ".") {
			return nil
		}

		annotationsPath = path
		// Ignore this error, since we only care if any annotations are returned.
		if annotations, err = readAnnotations(fs, path); err != nil {
			log.Debug(err)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	if len(annotations) == 0 {
		return nil, "", fmt.Errorf("metadata not found in %s", bundleRoot)
	}

	return annotations, annotationsPath, nil
}

// readAnnotations reads annotations from file(s) in bundleRoot and returns them as a map.
func readAnnotations(fs afero.Fs, annotationsPath string) (map[string]string, error) {
	// The annotations file is well-defined.
	b, err := afero.ReadFile(fs, annotationsPath)
	if err != nil {
		return nil, err
	}

	// Use the arbitrarily-indexed representation of the annotations file for forwards and backwards compatibility.
	annotations := struct {
		Annotations map[string]string `json:"annotations"`
	}{
		Annotations: make(map[string]string),
	}
	if err = yaml.Unmarshal(b, &annotations); err != nil {
		return nil, fmt.Errorf("error unmarshalling potential bundle metadata %s: %v", annotationsPath, err)
	}

	return annotations.Annotations, nil
}

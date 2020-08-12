package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// The bundle root-relative default annotations.yaml path.
var defaultRelPath = filepath.Join("metadata", "annotations.yaml")

// FindAnnotations walks bundleRoot searching for bundle metadata in an annotations.yaml file,
// and returns this metadata and its path if found. If one is not found, an error is returned.
func FindAnnotations(bundleRoot string) (AnnotationsFile, string, error) {
	return findAnnotations(afero.NewOsFs(), bundleRoot)
}

func findAnnotations(fs afero.Fs, bundleRoot string) (AnnotationsFile, string, error) {
	// Check the default path first, and return annotations if they were found or an error if that error
	// is not because the path does not exist (it exists or there was an unmarshaling error).
	annotationsPath := filepath.Join(bundleRoot, defaultRelPath)
	af, err := readAnnotations(fs, annotationsPath)
	if err != nil {
		// Ignore this error, since the annotations might be in some other file.
		log.Debug(err)
	} else if !af.IsEmpty() {
		return af, annotationsPath, nil
	}

	// Annotations are not at the default path, so search recursively.
	foundAnnotations := false
	annotationsPath = ""
	err = afero.Walk(fs, bundleRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip if annotations were already found, or path is a directory or hidden file.
		if foundAnnotations || info.IsDir() || strings.HasPrefix(path, ".") {
			return nil
		}

		annotationsPath = path
		if af, err = readAnnotations(fs, path); err != nil {
			// We don't want to stop early, so ignore this error and try other files.
			log.Debug(err)
		} else if !af.IsEmpty() {
			foundAnnotations = true
		}
		return nil
	})
	if err != nil {
		return af, "", err
	}

	if !foundAnnotations {
		return af, "", fmt.Errorf("metadata not found in %s", bundleRoot)
	}

	return af, annotationsPath, nil
}

// readAnnotations attempts to read annotations from annotationsPath.
func readAnnotations(fs afero.Fs, annotationsPath string) (af AnnotationsFile, err error) {
	b, err := afero.ReadFile(fs, annotationsPath)
	if err != nil {
		return af, err
	}
	if err = af.Unmarshal(b); err != nil {
		return af, fmt.Errorf("error unmarshaling potential bundle metadata %s: %v", annotationsPath, err)
	}
	return af, nil
}

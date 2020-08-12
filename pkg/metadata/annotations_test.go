package metadata

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Annotations", func() {
	Describe("FindAnnotations", func() {
		var (
			fs          afero.Fs
			err         error
			defaultPath = "/bundle/metadata/annotations.yaml"
		)

		Context("with valid annotations contents", func() {
			var (
				metadata      map[string]string
				path, expPath string
			)
			BeforeEach(func() {
				fs = afero.NewMemMapFs()
			})

			// Location
			It("finds registry metadata in the default location", func() {
				expPath = defaultPath
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})
			It("finds registry metadata in the a custom file name", func() {
				expPath = "/bundle/metadata/my-metadata.yaml"
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})
			It("finds registry metadata in a custom single-depth location", func() {
				expPath = "/bundle/my-dir/my-metadata.yaml"
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})
			It("finds registry metadata in a custom multi-depth location", func() {
				expPath = "/bundle/my-parent-dir/my-dir/annotations.yaml"
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})
			It("returns registry metadata from default path when metadata is also in another location", func() {
				expPath = defaultPath
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				writeMetadataHelper(fs, "/bundle/other-metadata/annotations.yaml", annotationsStringValidNoRegLabels)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})
			It("returns registry metadata from the first path, when metadata is also in another location", func() {
				expPath = "/bundle/custom1/annotations.yaml"
				writeMetadataHelper(fs, expPath, annotationsStringValidV1)
				writeMetadataHelper(fs, "/bundle/custom2/annotations.yaml", annotationsStringValidNoRegLabels)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidV1))
			})

			// Format
			It("finds non-registry metadata", func() {
				expPath = defaultPath
				writeMetadataHelper(fs, defaultPath, annotationsStringValidNoRegLabels)
				metadata, path, err = findAnnotations(fs, "/bundle")
				Expect(err).To(BeNil())
				Expect(path).To(Equal(expPath))
				Expect(metadata).To(BeEquivalentTo(annotationsValidNoRegLabels))
			})
		})

		Context("with invalid annotations contents", func() {
			BeforeEach(func() {
				fs = afero.NewMemMapFs()
			})

			It("returns a YAML error", func() {
				writeMetadataHelper(fs, defaultPath, annotationsStringInvalidBadIndent)
				_, _, err = findAnnotations(fs, "/bundle")
				// err should contain both of the following parts.
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("error unmarshalling potential bundle metadata %s: ", defaultPath)))
				Expect(err.Error()).To(ContainSubstring("yaml: line 2: found character that cannot start any token"))
			})
			It("returns an error for no metadata file (empty file)", func() {
				writeMetadataHelper(fs, defaultPath, annotationsStringInvalidEmpty)
				_, _, err = findAnnotations(fs, "/bundle")
				Expect(err).To(MatchError("metadata not found in /bundle"))
			})
			It("returns an error for no metadata file (invalid top-level key)", func() {
				writeMetadataHelper(fs, defaultPath, annotationsStringInvalidTopKey)
				_, _, err = findAnnotations(fs, "/bundle")
				Expect(err).To(MatchError("metadata not found in /bundle"))
			})
			It("returns an error for no labels in a metadata file", func() {
				writeMetadataHelper(fs, defaultPath, annotationsStringInvalidNoLabels)
				_, _, err = findAnnotations(fs, "/bundle")
				Expect(err).To(MatchError("metadata not found in /bundle"))
			})
		})
	})

})

func writeMetadataHelper(fs afero.Fs, path, contents string) {
	ExpectWithOffset(1, fs.MkdirAll(filepath.Dir(path), 0755)).Should(Succeed())
	ExpectWithOffset(1, afero.WriteFile(fs, path, []byte(contents), 0666)).Should(Succeed())
}

var annotationsValidV1 = map[string]string{
	"operators.operatorframework.io.bundle.mediatype.v1": "registry+v1",
	"operators.operatorframework.io.bundle.metadata.v1":  "metadata/",
	"foo": "bar",
}

const annotationsStringValidV1 = `annotations:
  operators.operatorframework.io.bundle.mediatype.v1: registry+v1
  operators.operatorframework.io.bundle.metadata.v1: metadata/
  foo: bar
`

var annotationsValidNoRegLabels = map[string]string{
	"foo": "bar",
	"baz": "buf",
}

const annotationsStringValidNoRegLabels = `annotations:
  foo: bar
  baz: buf
`

const annotationsStringInvalidBadIndent = `annotations:
	operators.operatorframework.io.bundle.mediatype.v1: registry+v1
`

const annotationsStringInvalidEmpty = ``

const annotationsStringInvalidNoLabels = `annotations:
`

const annotationsStringInvalidTopKey = `not-annotations:
  operators.operatorframework.io.bundle.mediatype.v1: registry+v1
  operators.operatorframework.io.bundle.metadata.v1: metadata/
  foo: bar
`

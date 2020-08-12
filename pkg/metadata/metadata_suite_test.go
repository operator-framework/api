package metadata

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetadata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metadata Suite")
}

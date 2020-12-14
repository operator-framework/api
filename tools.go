// +build tools

package tools

import (
	// Generate deepcopy and conversion.
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	// Manipulate YAML.
	_ "github.com/mikefarah/yq/v3"
	// Generate embedded files.
	_ "github.com/go-bindata/go-bindata/v3/go-bindata"
)

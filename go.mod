module github.com/operator-framework/api

go 1.15

require (
	github.com/blang/semver v3.5.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata/v3 v3.1.3
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mikefarah/yq/v2 v2.4.1
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.3.0
)

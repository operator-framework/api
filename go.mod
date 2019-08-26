module github.com/operator-framework/api

go 1.12

replace (
	// Pin openshift version to 4.2 (uses kube 1.14)
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20190717200738-0390d1e77d64+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190627172412-c44a8b61b9f4

	// Pin kube version to 1.14
	k8s.io/api => k8s.io/api v0.0.0-20190704095032-f4ca3d3bdf1d
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190704104557-6209bbe9f7a9
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190704094733-8f6ac2502e51
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190704101451-e5f5c6e528cd
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190521191137-11646d1007e0+incompatible
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190704094409-6c2a4329ac29
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190704100636-f0322db00a10
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190704101955-e796fd6d55e0
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.5-beta.0.0.20190708100021-7936da50c68f
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190926160646-a61144936680
	github.com/operator-framework/operator-registry v1.5.3
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.3.0
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v8.0.0+incompatible
)

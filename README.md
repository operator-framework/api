Api
--

Contains the API definitions used by [Operator Lifecycle Manager][olm] (OLM) and [Marketplace Operator][marketplace]

## Certificate of Origin
By contributing to this project you agree to the Developer Certificat of Origin (DCO).  This document was created by the Linux Kernel community and is a simple statement that you, as a contributor, have the legal right to make the contribution.  See the [DCO][DCO] file for details.

## `pkg/validation`: Operator Manifest Validation

`pkg/validation` exposes a convenient set of interfaces to validate Kubernetes object manifests, primarily for use in an Operator project.

The Validators are static checks (linters) that can scan the manifests and provide
with low-cost valuable results to ensure the quality of the package of distributions
(bundle or package formats) which will be distributed via OLM.

The validators implemented in this project aims to provide common validators
(which can be useful or required for any solution which will be distributed via [Operator Lifecycle Manager][olm]).
([More info](https://pkg.go.dev/github.com/operator-framework/api@master/pkg/validation))

Note that [Operator-SDK][sdk] leverage in this project. By using it you can
test your bundle against the spec criteria (Default Validators) by running:

```sh
$ operator-sdk bundle validate <bundle-path>
```

Also, [Operator-SDK][sdk] allows you check your bundles against the Optional Validators
provided by using the flag option `--select-optional` such as the following example:

```sh
$ operator-sdk bundle validate ./bundle --select-optional suite=operatorframework --optional-values=k8s-version=<k8s-version>
```

For further information see the [doc][sdk-command-doc].

### Example of usage:

Note that you can leverage in this project to call and indeed create your own validators.
Following an example.

```go
 import (
   ...
    apimanifests "github.com/operator-framework/api/pkg/manifests"
    apivalidation "github.com/operator-framework/api/pkg/validation"
    "github.com/operator-framework/api/pkg/validation/errors"
   ...
  )

 // Load the directory (which can be in packagemanifest or bundle format)
 bundle, err := apimanifests.GetBundleFromDir(path)
 if err != nil {
   ...
   return nil
 }

 // Call all default validators and the OperatorHubValidator
 validators := apivalidation.DefaultBundleValidators
 validators = validators.WithValidators(apivalidation.OperatorHubValidator)

 objs := bundle.ObjectsToValidate()

 results := validators.Validate(objs...)
 nonEmptyResults := []errors.ManifestResult{}

 for _, result := range results {
    if result.HasError() || result.HasWarn() {
        nonEmptyResults = append(nonEmptyResults, result)
    }
 }
 // return the results
 return nonEmptyResults
```

#### Passing optional key/values to the validators

Validators may accept pass optional key/values which will be used in the checks made.
These values are global and if the key/value pair provided is not used for 1 or more
validators called then, it is ignored.

The following example calls `AlphaDeprecatedAPIsValidator`, which allows us to inform
the K8s version intended to publish the OLM Bundle:

```go
	validators := apivalidation.DefaultBundleValidators
	validators = validators.WithValidators(apivalidation.OperatorHubValidator)
	validators = validators.WithValidators(apivalidation.ObjectValidator)
	validators = validators.WithValidators(apivalidation.AlphaDeprecatedAPIsValidator)
	validators = validators.WithValidators(apivalidation.GoodPracticesValidator)

	objs := auditBundle.Bundle.ObjectsToValidate()

	// Pass the --optional-values. e.g. --optional-values="k8s-version=1.22"
	// or --optional-values="image-path=bundle.Dockerfile"
	var optionalValues = map[string]string{
		"k8s-version":"1.22",
	}
	objs = append(objs, optionalValues)

	results := validators.Validate(objs...)
	nonEmptyResults := []errors.ManifestResult{}

	for _, result := range results {
		if result.HasError() || result.HasWarn() {
			nonEmptyResults = append(nonEmptyResults, result)
		}
	}
```

**How the optional key/values are informed via the CLI?**

By using [Operator-SDK][sdk] you can pass a list of key/values via the flag `--optional-values`, for example,
to validate that your manifests can work with a Kubernetes cluster of a particular version using the `k8s-version`:

```sh
$ operator-sdk bundle validate ./bundle --select-optional suite=operatorframework --optional-values=k8s-version=1.22
```

## API CLI Usage

You can install the `operator-verify` tool from source using:

`$ make install`

To verify your ClusterServiceVersion yaml,

`$ operator-verify manifests /path/to/filename.yaml`

[sdk]: https://github.com/operator-framework/operator-sdk
[olm]: https://github.com/operator-framework/operator-lifecycle-manager
[marketplace]: https://github.com/operator-framework/operator-marketplace
[bundle]: https://github.com/operator-framework/operator-registry/blob/v1.19.5/docs/design/operator-bundle.md
[sdk-command-doc]: https://master.sdk.operatorframework.io/docs/cli/operator-sdk_bundle_validate/

# api

Contains the API definitions used by [Operator Lifecycle Manager][olm] (OLM) and [Marketplace Operator][marketplace]

## `pkg/validation`: Operator Manifest Validation

`pkg/validation` exposes a convenient set of interfaces to validate Kubernetes object manifests, primarily for use in an Operator project.

[olm]:https://github.com/operator-framework/operator-lifecycle-manager
[marketplace]:https://github.com/operator-framework/operator-marketplace

## Usage

You can install the `operator-verify` tool from source using:

`$ make install`

To verify your ClusterServiceVersion yaml,

`$ operator-verify verify /path/to/filename.yaml`
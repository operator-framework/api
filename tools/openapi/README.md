`openapi/go.mod` pins the OpenAPI generator used by `make generate-openapi` to
`k8s.io/kube-openapi` from January 2026. The root `go.mod` requires the July
2026 version for the API library. Go selects one version of each module for a
module graph, so a tool directive in the root module would run the July
generator. That generator changes checked-in OpenAPI schemas, including
default values. Downgrading the root dependency to January would also change
the API library's Kubernetes dependency graph.

The Makefile runs the January generator with `go tool -modfile=tools/openapi/go.mod`.

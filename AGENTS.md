# AGENTS.md

This file provides AI agents with comprehensive context about the operator-framework/api repository to enable effective navigation, understanding, and contribution.

## Project Overview

This repository contains the API definitions and validation libraries used by [Operator Lifecycle Manager][olm] (OLMv0). It's a foundational library in the [Operator Framework](https://github.com/operator-framework) ecosystem.

### Core Capabilities
- **API Definitions**: Kubernetes Custom Resource Definitions (CRDs) for OLM resources
- **Manifest Validation**: Static validators for operator bundles and package manifests
- **CLI Tool**: `operator-verify` for manifest verification
- **Common Libraries**: Shared utilities for bundle and manifest manipulation

## Custom Resource Definitions (CRDs)

| Resource | API Group | Description |
|----------|-----------|-------------|
| **ClusterServiceVersion (CSV)** | operators.coreos.com/v1alpha1 | Defines operator metadata, installation strategy, permissions, and owned/required CRDs |
| **Subscription** | operators.coreos.com/v1alpha1 | Tracks operator updates from a catalog channel |
| **InstallPlan** | operators.coreos.com/v1alpha1 | Calculated list of resources to install/upgrade |
| **CatalogSource** | operators.coreos.com/v1alpha1 | Repository of operators and metadata |
| **OperatorGroup** | operators.coreos.com/v1 | Groups namespaces for operator installation scope |
| **OperatorCondition** | operators.coreos.com/v2 | Tracks operator health status and conditions |

## Directory Structure

```
api/
├── cmd/                          # Entry point binaries
│   └── operator-verify/          # CLI tool for manifest verification
│
├── pkg/                          # Core implementation
│   ├── operators/                # OLM API types
│   │   ├── v1alpha1/             # Core OLM types (CSV, Subscription, etc.)
│   │   ├── v1/                   # OperatorGroup, OperatorCondition v1
│   │   ├── v2/                   # OperatorCondition v2
│   │   └── reference/            # Image reference parsing
│   │
│   ├── validation/               # Operator manifest validators
│   │   ├── errors/               # Validation error types
│   │   ├── interfaces/           # Validator interfaces
│   │   └── internal/             # Validator implementations
│   │
│   ├── manifests/                # Bundle and manifest loaders
│   │
│   ├── constraints/              # Constraint and CEL validation
│   │
│   ├── lib/version/              # Version utilities
│   │
│   └── apis/scorecard/           # Scorecard configuration types
│
├── crds/                         # Generated CRD YAML files
│
└── hack/                         # Build scripts and tools
```

## Key Packages and Their Responsibilities

### API Types (`pkg/operators/`)

Defines all Kubernetes custom resources used by OLM:
- `v1alpha1/`: Core types (CSV, Subscription, InstallPlan, CatalogSource)
- `v1/`: OperatorGroup, Operator, OperatorCondition
- `v2/`: OperatorCondition v2
- `reference/`: Container image reference parsing utilities

**Key files**:
- `v1alpha1/clusterserviceversion_types.go` - CSV API definition
- `v1/operatorgroup_types.go` - OperatorGroup API definition

### Validation (`pkg/validation/`)

Static validators for operator bundles and manifests:
- **Default Validators**: Required checks for all operators
- **Optional Validators**: Community, OperatorHub, and best practice validators
- **Custom Validators**: Extensible validator interface

**Key files**:
- `validation.go` - Main validator orchestration
- `internal/bundle.go` - Bundle structure validation
- `internal/csv.go` - CSV validation rules
- `internal/operatorhub.go` - OperatorHub requirements

**Validator Types**:
- `BundleValidator` - Bundle format and structure
- `CSVValidator` - ClusterServiceVersion validation
- `CRDValidator` - CRD validation
- `OperatorHubValidator` - OperatorHub.io requirements
- `GoodPracticesValidator` - Best practices checks
- `AlphaDeprecatedAPIsValidator` - Deprecated API detection

### Manifests (`pkg/manifests/`)

Bundle and package manifest loaders:
- Bundle loading from directories
- PackageManifest parsing
- Metadata extraction

**Key files**:
- `bundle.go` - Bundle representation and loading
- `bundleloader.go` - Bundle loading logic

## Development Workflow

### Building operator-verify CLI
```bash
make install            # Build and install operator-verify CLI
```

### Testing
```bash
make test-unit          # Run unit tests
make test               # Run all tests
make TEST=<name> test-unit  # Run specific test
```

### Code Generation
```bash
make generate           # Generate deep-copy methods
make manifests          # Generate CRD manifests
make verify             # Verify generated code is up-to-date
```

### Code Quality
```bash
make format             # Format source code
make tidy               # Update and verify dependencies
```

## Validation Usage

### Using Default Validators

```go
import (
    apimanifests "github.com/operator-framework/api/pkg/manifests"
    apivalidation "github.com/operator-framework/api/pkg/validation"
)

// Load bundle
bundle, err := apimanifests.GetBundleFromDir(path)
if err != nil {
    return err
}

// Run default validators
validators := apivalidation.DefaultBundleValidators
results := validators.Validate(bundle.ObjectsToValidate()...)

// Check results
for _, result := range results {
    if result.HasError() {
        fmt.Printf("Error: %v\n", result)
    }
}
```

### Using Optional Validators

```go
// Add optional validators
validators := apivalidation.DefaultBundleValidators
validators = validators.WithValidators(apivalidation.OperatorHubValidator)
validators = validators.WithValidators(apivalidation.GoodPracticesValidator)

// Pass optional key/values
optionalValues := map[string]string{
    "k8s-version": "1.28",
}
objs := append(bundle.ObjectsToValidate(), optionalValues)

results := validators.Validate(objs...)
```

### CLI Usage

```bash
# Install operator-verify
make install

# Verify manifests
operator-verify manifests /path/to/manifest.yaml
```

## Code Generation

This repository uses controller-gen for code generation:

### Generated Code
- **Deep-copy methods**: Auto-generated for all API types (`zz_generated.deepcopy.go`)
- **CRD manifests**: Generated from Go type definitions in `crds/`
- **Embedded CRDs**: Go code embedding CRD YAML in `crds/zz_defs.go`

### Regenerating Code
```bash
make generate manifests  # Regenerate everything
make verify              # Verify nothing changed
```

**Important**: Never edit generated files directly - modify the source types and regenerate.

## Common Tasks for AI Agents

### Understanding Validation Flow
1. Bundle loaded via `pkg/manifests`
2. Validators instantiated from `pkg/validation`
3. Each validator checks specific aspects (CSV format, CRD structure, etc.)
4. Results aggregated with errors/warnings
5. Results returned to caller

### Adding a New Validator
1. Create new validator in `pkg/validation/internal/`
2. Implement `interfaces.Validator` interface
3. Add validator to appropriate suite in `pkg/validation/validation.go`
4. Write unit tests in `*_test.go`
5. Document validator behavior

### Modifying API Types
1. Edit type definition in `pkg/operators/v*/`
2. Update CRD markers (kubebuilder comments) if needed
3. Run `make generate manifests` to regenerate code
4. Run `make verify` to ensure clean state
5. Update tests if behavior changed

### Understanding CRD Structure
- All CRDs defined in `pkg/operators/v*/`
- Generated to `crds/*.yaml` via controller-gen
- Embedded in Go code at `crds/zz_defs.go`
- Used by OLM and other operator-framework components

## Important Dependencies

| Dependency | Purpose |
|------------|---------|
| k8s.io/api | Kubernetes core types |
| k8s.io/apimachinery | Kubernetes meta types |
| sigs.k8s.io/controller-tools | Code generation (controller-gen) |
| github.com/blang/semver/v4 | Semantic versioning |
| github.com/operator-framework/operator-registry | Registry integration |

## Navigation Tips

### Finding API Definitions
- Core OLM types: `pkg/operators/v1alpha1/`
- OperatorGroup: `pkg/operators/v1/operatorgroup_types.go`
- OperatorCondition: `pkg/operators/v2/operatorcondition_types.go`

### Finding Validators
- Validator implementations: `pkg/validation/internal/`
- Validator interfaces: `pkg/validation/interfaces/`
- Main orchestration: `pkg/validation/validation.go`

### Finding Utilities
- Bundle loading: `pkg/manifests/bundle.go`
- Image references: `pkg/operators/reference/`
- Version utilities: `pkg/lib/version/`

## Anti-Patterns to Avoid

1. **Don't modify generated code** - Edit source types and regenerate
2. **Don't skip `make verify`** - Always verify generated code is current
3. **Don't add breaking API changes** - This is a library used by multiple projects
4. **Don't add validators without tests** - All validators must be thoroughly tested
5. **Don't bypass validation interfaces** - Use the provided validator framework

## Resources and Links

- [OLM Documentation](https://olm.operatorframework.io/)
- [Operator SDK Bundle Validation](https://sdk.operatorframework.io/docs/cli/operator-sdk_bundle_validate/)
- [Operator Framework Community](https://github.com/operator-framework/community)
- [Validation Package Docs](https://pkg.go.dev/github.com/operator-framework/api@master/pkg/validation)

## Quick Reference

### Common Build Targets
```bash
make install             # Install operator-verify CLI
make test-unit           # Run unit tests
make generate manifests  # Generate code and CRDs
make verify              # Verify no uncommitted changes
make format tidy         # Format and tidy code
```

### Tool Management
Tools are managed via Makefile and installed to `./bin/`:
- `controller-gen` - Kubernetes code generator
- `yq` - YAML processor for CRD patching
- `kind` - Local Kubernetes clusters

## Contributing

See [DCO](DCO) for Developer Certificate of Origin requirements.

When contributing:
1. Run `make verify` before submitting PRs
2. Add tests for new validators or API changes
3. Update CRD generation if modifying types
4. Follow existing patterns for consistency

[olm]: https://github.com/operator-framework/operator-lifecycle-manager

# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
  BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
  Q =
else
  Q = @
endif

REPO = github.com/operator-framework/api
BUILD_PATH = $(REPO)/cmd/operator-verify
PKGS = $(shell go list ./... | grep -v /vendor/)

.PHONY: help
help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: install
install: ## Build & install operator-verify
	$(Q)go install \
		-gcflags "all=-trimpath=${GOPATH}" \
		-asmflags "all=-trimpath=${GOPATH}" \
		-ldflags " \
			-X '${REPO}/version.GitVersion=${VERSION}' \
			-X '${REPO}/version.GitCommit=${GIT_COMMIT}' \
		" \
		$(BUILD_PATH)

###
# Code management.
###
.PHONY: format tidy clean generate manifests

format: ## Format the source code
	$(Q)go fmt $(PKGS)

tidy: ## Update dependencies
	$(Q)go mod tidy
	$(Q)go mod verify

clean: ## Clean up the build artifacts
	$(Q)rm -rf build

generate: controller-gen  ## Generate code
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./...

manifests: yq controller-gen ## Generate manifests e.g. CRD, RBAC etc
	@# Create CRDs for new APIs
	$(CONTROLLER_GEN) crd:crdVersions=v1 output:crd:dir=./crds paths=./pkg/operators/...

	@# Update existing CRDs from type changes
	$(CONTROLLER_GEN) schemapatch:manifests=./crds output:dir=./crds paths=./pkg/operators/...

	@# Add missing defaults in embedded core API schemas
	$(YQ) --inplace '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.install.properties.spec.properties.deployments.items.properties.spec.properties.template.properties.spec.properties.containers.items.properties.ports.items.properties.protocol.default="TCP"' ./crds/operators.coreos.com_clusterserviceversions.yaml
	$(Q)$(YQ) --inplace '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.install.properties.spec.properties.deployments.items.properties.spec.properties.template.properties.spec.properties.initContainers.items.properties.ports.items.properties.protocol.default="TCP"' ./crds/operators.coreos.com_clusterserviceversions.yaml

	@# Preserve fields for embedded metadata fields
	$(Q)$(YQ) --inplace '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.install.properties.spec.properties.deployments.items.properties.spec.properties.template.properties.metadata.x-kubernetes-preserve-unknown-fields=true' ./crds/operators.coreos.com_clusterserviceversions.yaml

	@# Remove OperatorCondition.spec.overrides[*].lastTransitionTime requirement
	$(Q)$(YQ) --inplace 'del(.spec.versions[].schema.openAPIV3Schema.properties.spec.properties.overrides.items.required[] | select(. == "lastTransitionTime"))' ./crds/operators.coreos.com_operatorconditions.yaml

	@# Remove status subresource from the CRD manifests to ensure server-side apply works
	$(Q)for f in ./crds/*.yaml ; do $(YQ) --inplace 'del(.status)' $$f; done

	@# Update embedded CRD files.
	$(Q)go generate ./crds/...

# Static tests.
.PHONY: test test-unit verify

test: test-unit ## Run the tests

TEST_PKGS:=$(shell go list ./...)
test-unit: ## Run the unit tests
	$(Q)go test -count=1 -short ${TEST_PKGS}

verify: manifests generate format tidy
	git diff --exit-code

################
# Hack / Tools #
################

GO_INSTALL_OPTS ?= "-mod=mod"

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
YQ ?= $(LOCALBIN)/yq

## Tool Versions
CONTROLLER_TOOLS_VERSION ?= v0.14.0
YQ_VERSION ?= v4.28.1

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install $(GO_INSTALL_OPTS) sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: yq
yq: $(YQ) ## Download yq locally if necessary.
$(YQ): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install $(GO_INSTALL_OPTS) github.com/mikefarah/yq/v4@$(YQ_VERSION)


.PHONY: vendor codegen

vendor:
	go mod tidy
	go mod vendor

codegen: vendor
	./hack/update-codegen.sh
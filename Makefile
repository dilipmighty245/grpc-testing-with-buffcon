SHELL=/bin/bash

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif



.PHONY: test
test:
	set -o pipefail && go test -v -tags=unit -p=1 -count=1 -race -vet=off ./...
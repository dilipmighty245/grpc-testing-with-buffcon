SHELL=/bin/bash

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


.PHONY: integration-test
integration-test:
	set -o pipefail && go test -v -tags=integration -p=1 -count=1 -race -vet=off ./...

.PHONY: gen-proto
gen-proto:
	protoc ./proto/greeter/greeter.proto --go_out=plugins=grpc:.

.PHONY: gen-mocks
gen-mocks:
	${GOBIN}/mockgen -source=./proto/greeter/greeter.pb.go -destination=mocks/server_mock.go -package=mocks


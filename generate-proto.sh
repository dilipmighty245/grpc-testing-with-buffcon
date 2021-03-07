#!/bin/sh

#prerequisite:
# protobuf : https://developers.google.com/protocol-buffers/docs/gotutorial
# go: https://golang.org/doc/install
protoc ./proto/greeter/greeter.proto --go_out=plugins=grpc:.

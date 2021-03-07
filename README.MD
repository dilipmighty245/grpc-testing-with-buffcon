# Testing gRPC with buffcon

For the gRPC framework, we can run the gRPC server & client to test against the business logic. But spinning up a server from the test file could lead to unintended consequences that may require us to allocate a TCP port (parallel runs, multiple runs under the same CI server).

To solve this gRPC community has introduced a package called bufconn under gRPC’s testing package. bufconn is a package that provides a Listener object that implements net.Conn. We can substitute this listener in a gRPC server - allowing us to spin up a server that acts as a full-fledged server that can be used for testing that talks over an in-memory buffer instead of a real port.

As bufconn already comes with the grpc go module - which we already have installed, we don't need to install it explicitly.

## How to test
#### prerequisites:
```
Go : https://golang.org/doc/install
Git: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
Protobuf: https://developers.google.com/protocol-buffers/docs/gotutorial
```
#### Clone repo
```
git clone https://github.com/dilipmighty245/testing-grpc-with-buffcon.git
cd testing-grpc-with-buffcon
````
#### Generate proto
```
make gen-proto
```
#### Test server
```
make test

set -o pipefail && go test -v -tags=unit -p=1 -count=1 -race -vet=off ./...
?   	github.com/dilipmighty/testing-grpc-with-bufconn/proto/greeter	[no test files]
=== RUN   TestSayHello
2021/03/07 13:33:49 Received: gRPC
--- PASS: TestSayHello (0.01s)
PASS
ok  	github.com/dilipmighty/testing-grpc-with-bufconn/server	0.510s
```




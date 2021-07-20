# Testing gRPC rpcs

- Testing gRPC with gomock
- Testing gRPC with custom mock
- Testing gRPC on port
- Testing gRPC with buffcon


## How to test
#### Prerequisites:
```
Go : https://golang.org/doc/install
Protobuf: https://developers.google.com/protocol-buffers/docs/gotutorial
```
#### Clone repo
```
git clone https://github.com/dilipmighty245/testing-grpc-with-buffcon.git
cd testing-grpc-with-buffcon
````
#### For Generating proto stub from file
```
make gen-proto 
```
#### Test server
```
$ make integration-test

set -o pipefail && go test -v -tags=integration -p=1 -count=1 -race -vet=off ./...
Running Suite: Server BBD Test Suite
====================================
Random Seed: 1622099839
Will run 1 of 1 specs

2021/05/27 12:47:19 Received: gRPC
•2021/05/27 12:47:19 context is cancelled, proceeding for graceful shutdown

Ran 1 of 1 Specs in 0.007 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestClient (0.01s)
PASS
ok  	github.com/dilipmighty/testing-grpc-with-bufconn/server	0.440s


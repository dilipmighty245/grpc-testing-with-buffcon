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
?       github.com/dilipmighty/testing-grpc/mocks       [no test files]
?       github.com/dilipmighty/testing-grpc/proto/greeter       [no test files]
=== RUN   TestServer
Running Suite: Server BBD Test Suite
====================================
Random Seed: 1626781419
Will run 4 of 4 specs

•starting server at :8080
2021/07/20 17:13:39 Received: gRPC
•2021/07/20 17:13:39 Received: gRPC
••context is cancelled, proceeding for graceful shutdown

Ran 4 of 4 Specs in 0.015 seconds
SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestServer (0.03s)
PASS
ok      github.com/dilipmighty/testing-grpc/server      1.039s




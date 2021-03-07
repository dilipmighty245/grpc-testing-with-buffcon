# Testing gRPC with buffcon

For the gRPC framework, we can run the gRPC server & client to test against the business logic. But spinning up a server from the test file could lead to unintended consequences that may require us to allocate a TCP port (parallel runs, multiple runs under the same CI server).

To solve this gRPC community has introduced a package called bufconn under gRPC’s testing package. bufconn is a package that provides a Listener object that implements net.Conn. We can substitute this listener in a gRPC server - allowing us to spin up a server that acts as a full-fledged server that can be used for testing that talks over an in-memory buffer instead of a real port.

As bufconn already comes with the grpc go module - which we already have installed, we don't need to install it explicitly.

// +build unit

package main

import (
	"context"
	pb "github.com/dilipmighty/grpc-testing-with-bufconn/proto/greeter"
	"net"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis    *bufconn.Listener
	wg     sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	lis = bufconn.Listen(bufSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		run(ctx, lis)
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
	if err != nil {
		t.Fatal(err)
	}

	if resp.GetMessage() != "Hello gRPC" {
		t.Fatal("hello reply must be 'Hello test'")
	}
	cancel()
	wg.Wait()
}

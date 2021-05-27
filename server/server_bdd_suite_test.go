// +build integration

package main

import (
	"context"
	"fmt"
	pb "github.com/dilipmighty/testing-grpc-with-bufconn/proto/greeter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server BBD Test Suite")
}

var (
	lis    *bufconn.Listener
	wg     sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context
	s      *grpc.Server
	conn   *grpc.ClientConn
	err    error
)

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.Background())
	bufSize := 1024 * 1024
	lis = bufconn.Listen(bufSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer lis.Close()
		run(ctx, lis) // blocking call
	}()
	s = grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
})
var _ = Describe("Test Server with Port", func() {
	Context("Given a gRPC server with buffered connection, when hello api is called", func() {
		When("Hello api is called", func() {
			BeforeEach(func() {
				conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
					return lis.Dial()
				}), grpc.WithInsecure(), grpc.WithBlock())
				Expect(err).ToNot(HaveOccurred())
				Expect(conn).ToNot(BeNil())
			})

			It("should return valid response from server", func() {
				client := pb.NewGreeterClient(conn)
				resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.GetMessage()).To(Equal("Hello gRPC"))
			})
		})
	})
})

var _ = AfterSuite(func() {
	fmt.Fprintf(GinkgoWriter, " <-- Gracefully stopping server\n")
	const maxShutdownWait = 5 * time.Second
	doneCh := make(chan bool)
	go func() {
		defer conn.Close()
		wg.Wait()
		doneCh <- true
	}()

	cancel()
	s.GracefulStop()
	log.Printf("context is cancelled, proceeding for graceful shutdown")
	select {
	case <-time.After(maxShutdownWait):
		Fail("Timed Out during graceful shutdown")
	case <-doneCh:
		fmt.Fprintf(GinkgoWriter, " <-- Graceful shutdown\n")
	}
})

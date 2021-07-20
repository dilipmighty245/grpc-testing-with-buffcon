// +build integration

package main

import (
	"context"
	"fmt"
	"github.com/dilipmighty/testing-grpc/mocks"
	pb "github.com/dilipmighty/testing-grpc/proto/greeter"
	"github.com/golang/mock/gomock"
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

var testingT *testing.T

func TestServer(t *testing.T) {
	testingT = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client BBD Test Suite")
}

var (
	wg   sync.WaitGroup
	conn *grpc.ClientConn
	err  error
)

type server struct{}

func (s server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", request.Name)
	return &pb.HelloReply{Message: "Hello " + request.Name}, nil
}

var _ = Describe("Test Server with Port", func() {
	var s *grpc.Server

	Context("Given a gRPC server with tcp connection", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())

			s = grpc.NewServer()
			pb.RegisterGreeterServer(s, server{})

			addr := ":8080"
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			fmt.Println("starting server at " + addr)

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := run(ctx, lis); err != nil {
					fmt.Println("Server exited with error", err)
				}
			}()

			conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
			Expect(err).ToNot(HaveOccurred())
			Expect(conn).ToNot(BeNil())
		})

		It("should return valid response from server", func() {
			fakeClient := pb.NewGreeterClient(conn)
			resp, err := fakeClient.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetMessage()).To(Equal("Hello gRPC"))

		})
		AfterEach(func() {
			cancel()
			// s.GracefulStop()
		})
	})
})

const bufSize = 1024 * 1024

var _ = Describe("Test Server with buff conn", func() {

	Context("Given a gRPC server with buffered connection", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())

			lis := bufconn.Listen(bufSize)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := run(ctx, lis); err != nil {
					fmt.Println("Server exited with error", err)
				}
			}()

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
		AfterEach(func() {
			cancel()
			// s.GracefulStop()
		})
	})
})

var _ = Describe("Test Server with gomock", func() {

	Context("Given a gRPC server with tcp connection", func() {
		var mockCtrl *gomock.Controller
		var mockClient *mocks.MockGreeterClient
		req := &pb.HelloRequest{Name: ""}
		BeforeEach(func() {
			mockCtrl = gomock.NewController(testingT)
			mockClient = mocks.NewMockGreeterClient(mockCtrl)

			resp := &pb.HelloReply{
				Message: "Hello gRPC",
			}
			mockClient.EXPECT().SayHello(gomock.Any(), gomock.Any()).Return(resp, nil) // set the expectation
		})

		It("should return valid response", func() {
			resp, err := mockClient.SayHello(context.Background(), req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetMessage()).To(Equal("Hello gRPC"))
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})
	})
})

type FakeClient struct {
}

func NewClient() *FakeClient {
	return &FakeClient{}
}
func (c *FakeClient) SayHello(ctx context.Context, in *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello gRPC",
	}, nil
}

var _ pb.GreeterClient = (*FakeClient)(nil)

var _ = Describe("Test Server with custom mock", func() {

	Context("Given a gRPC server with tcp connection", func() {
		It("should return valid response from server", func() {
			client := NewClient()
			resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "gRPC"})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetMessage()).To(Equal("Hello gRPC"))
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

	fmt.Println("context is cancelled, proceeding for graceful shutdown")
	select {
	case <-time.After(maxShutdownWait):
		Fail("Timed Out during graceful shutdown")
	case <-doneCh:
		fmt.Fprintf(GinkgoWriter, " <-- Graceful shutdown\n")
	}
})

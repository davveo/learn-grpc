package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/proto/helloworld"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "the port to serve on")

// hwServer is used to implement helloworld.GreeterServer.
type hwServer struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hwServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

type ecServer struct {
	echo.UnimplementedEchoServer
}

func (s *ecServer) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: req.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()

	// Register Greeter on the server.
	helloworld.RegisterGreeterServer(s, &hwServer{})

	// Register RouteGuide on the same server.
	echo.RegisterEchoServer(s, &ecServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

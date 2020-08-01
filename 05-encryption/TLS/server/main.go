package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var port = flag.Int("port", 50051, "the port to serve on")

type Service struct {
	echo.UnimplementedEchoServer
}

func (s *Service) Simple(ctx context.Context, request *echo.Request) (*echo.Response, error) {
	log.Println("收到来自客户端的请求...")
	return &echo.Response{Message: request.Message}, nil
}

func (s *Service) ClientStream(server echo.Echo_ClientStreamServer) error {
	panic("implement me")
}

func (s *Service) ServerStream(request *echo.Request, server echo.Echo_ServerStreamServer) error {
	panic("implement me")
}

func (s *Service) DoubleStream(server echo.Echo_DoubleStreamServer) error {
	panic("implement me")
}

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(
		testdata.Path("server1.pem"),
		testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	echo.RegisterEchoServer(s, &Service{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

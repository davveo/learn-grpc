package main

import (
	"context"
	"flag"
	"fmt"
	echo "github.com/davveo/learn-grpc/pb"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "the port to serve on")

type service struct {
	pb echo.UnimplementedEchoServer
}

func (s *service) ClientStream(stream echo.Echo_ClientStreamServer) error {
	fmt.Printf("--- ClientStreamingEcho ---\n")

	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&echo.Response{Message: message})
		}
		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

func (s *service) ServerStream(request *echo.Request, server echo.Echo_ServerStreamServer) error {
	panic("implement me")
}

func (s *service) DoubleStream(server echo.Echo_DoubleStreamServer) error {
	panic("implement me")
}

func (s *service) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: req.Message}, nil
}

func main()  {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &service{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	utils.CheckError(err)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

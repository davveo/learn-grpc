package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 50051, "the port to serve on")
	streamingCount = 10
)

type service struct {
	pb echo.UnimplementedEchoServer
}

func (s *service) ClientStream(stream echo.Echo_ClientStreamServer) error {
	panic("implement me")
}

func (s *service) ServerStream(request *echo.Request, stream echo.Echo_ServerStreamServer) error {
	fmt.Printf("--- ServerStreamingEcho ---\n")

	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", request.Message)
		err := stream.Send(&echo.Response{Message: request.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DoubleStream(server echo.Echo_DoubleStreamServer) error {
	panic("implement me")
}

func (s *service) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: req.Message}, nil
}

func main() {
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
